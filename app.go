package eeaao_codegen

import (
	"encoding/json"
	"github.com/palindrom615/eeaao-codegen/codelet"
	json2 "go.starlark.net/lib/json"
	"go.starlark.net/lib/math"
	"go.starlark.net/lib/time"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
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

// Render renders the templates.
// internally, it just runs `render.star` file in the CodeletDir
func (a *App) Render() {
	err := os.MkdirAll(a.OutDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating output directory: %v\n", err)
	}
	// run `render.star` in the CodeletDir
	starlarkThread := &starlark.Thread{
		Name: "main",
	}
	predefined := starlark.StringDict{
		"eeaao": codelet.ToStarlarkModule(a),
		"json":  json2.Module,
		//"proto": proto.Module,
		"math": math.Module,
		"time": time.Module,
	}
	globals, err := starlark.ExecFileOptions(
		&syntax.FileOptions{},
		starlarkThread,
		filepath.Join(a.CodeletDir, "render.star"),
		nil,
		predefined,
	)
	if err != nil {
		log.Printf("Error running render.star. globals: %v\n%v\n", globals, err)
	}
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
