package eeaao_codegen

import (
	"encoding/json"
	"github.com/Masterminds/sprig"
	"github.com/palindrom615/eeaao-codegen/plugin"
	"gopkg.in/yaml.v3"
	"log"
	"maps"
	"os"
	"path/filepath"
	"text/template"
)

type App struct {
	SpecDir    string
	OutDir     string
	CodeletDir string
	Conf       map[string]any
}

func NewApp(specDir string, outDir string, codeletDir string, configFile string) *App {
	conf := readConf(configFile)
	return &App{
		SpecDir:    specDir,
		OutDir:     outDir,
		CodeletDir: codeletDir,
		Conf:       conf,
	}
}

func (a *App) renderFile(filePath string, templatePath string, data any) string {
	dst := filepath.Join(a.OutDir, filePath)
	dstFile, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	// TODO deny access to files outside of the codelet directory
	tmplStr, err := os.ReadFile(filepath.Join(a.CodeletDir, "templates", templatePath))
	tmpl := template.New(filepath.Base(filepath.Join(a.CodeletDir, "templates", templatePath)))
	tmpl.Parse(string(tmplStr))
	tmpl.Execute(dstFile, data)

	return dst
}

func (a *App) loadSpecsGlob(pluginName string, glob string) (res []plugin.SpecData) {
	p := plugin.GetPlugin(pluginName)
	matches, err := filepath.Glob(filepath.Join(a.SpecDir, glob))
	if err != nil {
		return nil
	}
	for _, match := range matches {
		if doc := p.LoadSpecFile(match); doc != nil {
			res = append(res, doc)
		}
	}
	return res
}

func readConf(configFile string) map[string]any {
	config := make(map[string]any)
	if configFile == "" {
		return config
	}
	configData, err := os.ReadFile(configFile)
	if err != nil {
		log.Printf("config file not found: %s", configFile)

	}
	ext := filepath.Ext(configFile)
	if ext == ".json" {
		err = json.Unmarshal(configData, &config)
		if err != nil {
			panic(err)
		}
	} else if ext == ".yaml" || ext == ".yml" {
		err = yaml.Unmarshal(configData, &config)
		if err != nil {
			panic(err)
		}
	}
	return config
}

func (c *App) Render() string {
	// Read JSON specFile

	os.Mkdir(c.OutDir, os.ModePerm)

	// render `render.tmpl` with spec data
	tmplData, err := os.ReadFile(filepath.Join(c.CodeletDir, "render.tmpl"))
	tmpl := template.New("render.tmpl")
	funcmap := template.FuncMap{
		"loadSpecsGlob": c.loadSpecsGlob,
		"renderFile":    c.renderFile,
		"withConfig":    func() map[string]any { return c.Conf },
	}
	maps.Copy(funcmap, sprig.FuncMap())
	tmpl.Funcs(funcmap)

	tmpl, err = tmpl.Parse(string(tmplData))
	if err != nil {
		log.Fatal(err)
		return ""
	}

	// Create output file
	outFilePath := filepath.Join(c.OutDir, "render")
	outFile, err := os.Create(outFilePath)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer outFile.Close()

	// Render template with spec data
	err = tmpl.Execute(outFile, struct{}{})
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return outFilePath
}
