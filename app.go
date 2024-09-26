package eeaao_codegen

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"log"
	"maps"
	"os"
	"path/filepath"
	"text/template"
)

type App struct {
	specDir        string
	OutDir         string
	CodeletDir     string
	Values         map[string]any
	tmpl           *template.Template
	starlarkRunner *starlarkRunner
}

// NewApp creates a new App instance
// specDir: directory for specifications
// outDir: directory for output
// codeletDir: directory for codelet
// valuesFile: file path of external values file. if empty, it will be ignored.
func NewApp(specDir string, outDir string, codeletDir string, valuesFile string) *App {
	a := &App{
		specDir:    specDir,
		OutDir:     outDir,
		CodeletDir: codeletDir,
	}
	a.loadValues(valuesFile)

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

// loadValues loads App.Values from valuesFile.
// it first tries to load default values from CodeletDir/values.yaml.
// then it tries to load external values from valuesFile and override the default values.
// any error will be logged and ignored.
func (a *App) loadValues(valuesFile string) {
	a.Values = make(map[string]any)
	defaultValueFilePath := filepath.Join(a.CodeletDir, "values.yaml")
	defaultValueFile, err := os.Open(defaultValueFilePath)
	if err != nil {
		log.Printf("Error opening default values file %s: %v\n", defaultValueFilePath, err)
	}
	err = yaml.NewDecoder(defaultValueFile).Decode(&a.Values)
	if err != nil {
		log.Printf("Error parsing default values file %s: %v\n", defaultValueFilePath, err)
	}
	if valuesFile == "" {
		return
	}
	c, err := os.Open(valuesFile)
	if err != nil {
		log.Printf("Error opening values file: %s\n%v\n", valuesFile, err)
		return
	}

	externalValues := make(map[string]any)
	ext := filepath.Ext(valuesFile)
	if ext == ".json" {
		err = json.NewDecoder(c).Decode(&externalValues)
	} else if ext == ".yaml" || ext == ".yml" {
		err = yaml.NewDecoder(c).Decode(&externalValues)
	} else {
		log.Printf("Unknown values file type: %s\n", valuesFile)
	}
	if err != nil {
		log.Printf("Error parsing values file: %s\n%v\n", valuesFile, err)
	}
	maps.Copy(a.Values, externalValues)
}

// Render renders the templates.
// internally, it just runs `render.star` file in the CodeletDir
func (a *App) Render() error {
	_, err := a.starlarkRunner.Render()
	return err
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
