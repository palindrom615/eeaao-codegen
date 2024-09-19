package eeaao_codegen

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
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
// codeletDir: directory for codelet
// configFile: config file. if empty, no config is loaded.
func NewApp(specDir string, outDir string, codeletDir string, configFile string) *App {
	conf := readConf(configFile)

	a := &App{
		specDir:    specDir,
		OutDir:     outDir,
		CodeletDir: codeletDir,
		Conf:       conf,
	}
	a.populateTemplate()
	runner, err := newStarlarkRunner(codeletDir, ToStarlarkModule(a))
	if err != nil {
		log.Fatalf("Error creating starlark runner: %v\n", err)
	}
	a.starlarkRunner = runner
	err = os.MkdirAll(a.OutDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating output directory: %v\n", err)
	}
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
	} else if ext == ".yaml" || ext == ".yml" {
		err = yaml.Unmarshal(configData, &config)
	}
	if err != nil {
		log.Printf("Error parsing config file: %s\n%v\n", configFile, err)
	}
	return config
}

// Render renders the templates.
// internally, it just runs `render.star` file in the CodeletDir
func (a *App) Render() {
	a.starlarkRunner.Render()
}

// RunShell starts a REPL shell for testing
func (a *App) RunShell() {
	a.starlarkRunner.RunShell()
}

func (a *App) populateTemplate() {
	a.tmpl = template.New("root")
	a.tmpl.Funcs(ToTemplateFuncmap(a))
	tmplDir := filepath.Join(a.CodeletDir, "templates")
	if stat, err := os.Stat(tmplDir); err != nil || !stat.IsDir() {
		log.Printf("Failed to find templates directory [%s]\n", tmplDir)
		return
	}

	filepath.Walk(tmplDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		tmplName, _ := filepath.Rel(tmplDir, path)
		tmplText, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Failed reading template file [%s]. skipped\n%v\n", path, err)
			return nil
		}
		_, err = a.tmpl.New(tmplName).Parse(string(tmplText))
		if err != nil {
			log.Printf("Failed parsing template file [%s]; skipped\n%v\n", path, err)
			return nil
		}
		return nil
	})
}
