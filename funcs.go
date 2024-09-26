package eeaao_codegen

import (
	"bytes"
	"encoding/json"
	"github.com/palindrom615/eeaao-codegen/plugin"
	"log"
	"os"
	"path/filepath"
)

func (a *App) RenderFile(filePath string, templatePath string, data any) (dst string, err error) {
	if !filepath.IsLocal(filePath) {
		log.Printf("invalid filePath: %s", filePath)
		return "", nil
	}
	if !filepath.IsLocal(templatePath) {
		log.Printf("invalid templatePath: %s", templatePath)
		return "", nil
	}
	dst = filepath.Join(a.OutDir, filePath)
	os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	dstFile, err := os.Create(dst)
	if err != nil {
		log.Printf("Error creating file '%s': %v\n", dst, err)
		return "", err
	}
	err = a.tmpl.ExecuteTemplate(dstFile, templatePath, data)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

func (a *App) LoadSpecsGlob(pluginName string, glob string) (map[string]string, error) {
	p := plugin.GetPlugin(pluginName)
	matches, err := filepath.Glob(filepath.Join(a.specDir, glob))
	if err != nil {
		return nil, err
	}
	res := make(map[string]string)
	for _, match := range matches {
		doc, err := p.LoadSpecFile(match)
		if err != nil {
			log.Printf("Error loading spec file '%s': %v\n", match, err)
			continue
		}
		docStr, err := json.Marshal(doc)
		if err != nil {
			log.Printf("Error marshaling spec file '%s': %v\n", match, err)
			continue
		}
		p, _ := filepath.Rel(a.specDir, match)
		res[p] = string(docStr)
	}
	return res, nil
}

func (a *App) LoadValues() map[string]any {
	return a.Values
}

func (a *App) Include(templatePath string, data interface{}) (string, error) {
	buf := bytes.NewBuffer(nil)
	if err := a.tmpl.ExecuteTemplate(buf, templatePath, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
