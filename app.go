package eeaao_codegen

import (
	"encoding/json"
	"github.com/palindrom615/eeaao-codegen/codelet"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type App struct {
	specDir        string
	OutDir         string
	CodeletDir     string
	Conf           map[string]any
	tmpl           *template.Template
	starlarkRunner *starlarkRunner
}

// NewApp creates a new App instance
// specDir: directory for specifications
// outDir: directory for output
// codeletDir: directory for templates
// configFile: config file. if empty, config is ""
func NewApp(specDir string, outDir string, codeletDir string, configFile string) *App {
	conf := readConf(configFile)
	a := &App{
		specDir:        specDir,
		OutDir:         outDir,
		CodeletDir:     codeletDir,
		Conf:           conf,
		starlarkRunner: newStarlarkRunner(),
	}
	a.populateTemplate()
	err := os.MkdirAll(a.OutDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating output directory: %v\n", err)
	}
	a.starlarkRunner.addModule("eeaao_codegen", codelet.ToStarlarkModule(a))
	return a
}

func readConf(configFile string) map[string]any {
	config := make(map[string]any)
	if configFile == "" {
		return config
	}
	configData, err := os.ReadFile(configFile)
	if err != nil {
		log.Printf("config file not found: %s\n%v\n", configFile, err)
		return config
	}
	ext := filepath.Ext(configFile)
	if ext == ".json" {
		err = json.Unmarshal(configData, &config)
		if err != nil {
			log.Printf("Error parsing config file: %s\n%v\n", configFile, err)
		}
	} else if ext == ".yaml" || ext == ".yml" {
		err = yaml.Unmarshal(configData, &config)
		if err != nil {
			log.Printf("Error parsing config file: %s\n%v\n", configFile, err)
		}
	}
	return config
}

// Render renders the templates.
// internally, it just runs `render.star` file in the CodeletDir
func (a *App) Render() {
	globals, err := a.starlarkRunner.runFile(filepath.Join(a.CodeletDir, "render.star"))
	if err != nil {
		log.Printf("Error running render.star. globals: %v\n%v\n", globals, err)
	}
}

// RunShell starts a REPL shell for testing
func (a *App) RunShell() {
	a.starlarkRunner.runShell()
}

func (a *App) populateTemplate() {
	a.tmpl = template.New("root")
	a.tmpl.Funcs(codelet.ToTemplateFuncmap(a))
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
