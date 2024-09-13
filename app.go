package eeaao_codegen

import (
	"encoding/json"
	"github.com/palindrom615/eeaao-codegen/plugin"
	"gopkg.in/yaml.v3"
	"log"
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

func (a *App) LoadSpecsGlob(pluginName string, glob string) (res []plugin.SpecData) {
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
