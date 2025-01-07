package eeaao_codegen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/palindrom615/eeaao-codegen/plugin"
	"gopkg.in/yaml.v3"
	"log"
	"maps"
	"os"
	"path/filepath"
	"text/template"
)

type App struct {
	OutDir         string
	CodeletDir     string
	Values         map[string]any
	tmpl           *template.Template
	starlarkRunner *starlarkRunner
	plugins        plugin.Plugins
}

// NewApp creates a new App instance
// outDir: directory for output
// codeletDir: directory for codelet
// valuesFile: file path of external values file. if empty, it will be ignored.
func NewApp(outDir string, codeletDir string, valuesFile string) *App {
	a := &App{
		OutDir:     outDir,
		CodeletDir: codeletDir,
		plugins:    plugin.NewPlugins(),
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

// RenderFile renders a file with the given template and data
// filePath: the file path to render. The path is relative to the output directory.
// templatePath: the template path. The path is relative to the ${codeletdir}/templates directory.
// data: the data to render
// returns the destination file path.
func (a *App) RenderFile(filePath string, templatePath string, data any) (dst string, err error) {
	if !filepath.IsLocal(filePath) {
		return "", fmt.Errorf("invalid filePath: %s", filePath)
	}
	if !filepath.IsLocal(templatePath) {
		return "", fmt.Errorf("invalid templatePath: %s", templatePath)
	}
	dst = filepath.Join(a.OutDir, filePath)
	os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	dstFile, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("error creating file '%s': %w", dst, err)
	}
	err = a.tmpl.ExecuteTemplate(dstFile, templatePath, data)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

// LoadValues returns the values data from codelet's default values.yaml file and given values file.
func (a *App) LoadValues() map[string]any {
	return a.Values
}

// GetPlugin returns the plugin with the given name.
func (a *App) GetPlugin(pluginName string) plugin.Plugin {
	return a.plugins.GetPlugin(pluginName)
}

// Include renders a template with the given data.
//
// Drop-in replacement for template pipeline, but with a string return value so that it can be treated as a string in the template.
//
// Inspired by [helm include function](https://helm.sh/docs/chart_template_guide/named_templates/#the-include-function)
func (a *App) Include(templatePath string, data interface{}) (string, error) {
	buf := bytes.NewBuffer(nil)
	if err := a.tmpl.ExecuteTemplate(buf, templatePath, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
