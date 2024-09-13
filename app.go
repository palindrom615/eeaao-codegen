package eeaao_codegen

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type App struct {
	specDir    string
	OutDir     string
	CodeletDir string
	Conf       map[string]any
	tmpl       *template.Template
}

// NewApp creates a new App instance
// specDir: directory for specifications
// outDir: directory for output
// codeletDir: directory for templates
// configFile: config file. if empty, config is ""
func NewApp(specDir string, outDir string, codeletDir string, configFile string) *App {
	conf := readConf(configFile)
	app := &App{
		specDir:    specDir,
		OutDir:     outDir,
		CodeletDir: codeletDir,
		Conf:       conf,
	}
	app.populateTemplate()
	return app
}

func readConf(configFile string) map[string]any {
	config := make(map[string]any)
	if configFile == "" {
		return config
	}
	configData, err := os.ReadFile(configFile)
	if err != nil {
		log.Printf("config file not found: %s\n%v\n", configFile, err)
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

// Render renders the templates
func (a *App) Render() string {
	err := os.MkdirAll(a.OutDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating output directory: %v\n", err)
	}

	// render `render.tmpl` with spec data
	tmplData, err := os.ReadFile(filepath.Join(a.CodeletDir, "render.tmpl"))
	tmpl := a.tmpl.New("../render.tmpl")

	tmpl, err = tmpl.Parse(string(tmplData))
	if err != nil {
		log.Fatalf("Error parsing render.tmpl: %v\n", err)
	}

	// Create output file
	outFilePath := filepath.Join(a.OutDir, "render")
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

func (a *App) populateTemplate() {
	a.tmpl = template.New("root")
	a.tmpl.Funcs(a.makeFuncmap())
	tmplDir := filepath.Join(a.CodeletDir, "templates")

	filepath.Walk(tmplDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		tmplName, found := strings.CutPrefix(path, tmplDir+"/")
		if !found {
			log.Fatalf("Error: %s is not in %s\n", path, tmplDir)
		}
		tmplText, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Error reading template file: %s; ignored\n%v\n", path, err)
			return nil
		}
		_, err = a.tmpl.New(tmplName).Parse(string(tmplText))
		if err != nil {
			log.Printf("Error parsing template file: %s; ignored\n%v\n", path, err)
			return nil
		}
		return nil
	})
}
