package eeaao_codegen

import (
	"bytes"
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

func (a *App) LoadValues() map[string]any {
	return a.Values
}

func (a *App) GetPlugin(pluginName string) plugin.Plugin {
	return a.plugins.GetPlugin(pluginName)
}

func (a *App) Include(templatePath string, data interface{}) (string, error) {
	buf := bytes.NewBuffer(nil)
	if err := a.tmpl.ExecuteTemplate(buf, templatePath, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
