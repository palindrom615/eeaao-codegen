package eeaao_codegen

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Template struct {
	*template.Template
	tmplDir string
}

// NewTemplate creates a new Template instance.
// It initializes the template with the given directory and additional functions.
// The additional functions are
//   - sprig.FuncMap
//   - funcMap given as an argument.
//   - include: Template.Include method
//
// tmplDir: directory for templates
// funcMap: additional functions to be added to the template.
func NewTemplate(tmplDir string, funcMap template.FuncMap) (t *Template) {
	t = &Template{template.New("root"), tmplDir}
	t.Funcs(sprig.FuncMap())
	t.Funcs(funcMap)
	t.AddTemplateFunc("include", t.Include)
	return t
}

// AddTemplateFunc adds a function to the template.
func (t *Template) AddTemplateFunc(name string, f interface{}) {
	t.Funcs(template.FuncMap{name: f})
}

// ExecuteTemplate executes the template specified by the tmplPath, which is relative to the template directory.
// If the template is not found in the template directory, it reads the template file from the template directory
// and parses it.
func (t *Template) ExecuteTemplate(wr io.Writer, tmplPath string, data any) error {
	tmpl := t.Lookup(tmplPath)
	if tmpl == nil {
		tmplText, err := os.ReadFile(filepath.Join(t.tmplDir, tmplPath))
		if err != nil {
			return fmt.Errorf("error reading template file [%s]: %w", tmplPath, err)
		}
		_, err = t.New(tmplPath).Parse(string(tmplText))
		if err != nil {
			return fmt.Errorf("error parsing template file [%s]: %w", tmplPath, err)
		}
	}
	return t.Template.ExecuteTemplate(wr, tmplPath, data)
}

// Include renders a template with the given data.
//
// Drop-in replacement for template pipeline, but with a string return value so that it can be treated as a string in the template.
//
// Inspired by [helm include function](https://helm.sh/docs/chart_template_guide/named_templates/#the-include-function)
func (t *Template) Include(tmplPath string, data any) (string, error) {
	var sb strings.Builder
	err := t.ExecuteTemplate(&sb, tmplPath, data)
	return sb.String(), err
}
