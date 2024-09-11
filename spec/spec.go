package spec

import (
	"github.com/palindrom615/eeaao-codegen/plugin"
	"path/filepath"
)

type Spec any

type SpecDir struct {
	Dir string
}

func (s SpecDir) LoadSpecsGlob(pluginName string, glob string) []Spec {
	p := plugin.Plugins[pluginName]
	matches, err := filepath.Glob(filepath.Join(s.Dir, glob))
	if err != nil {
		return nil
	}
	var res []Spec
	for _, match := range matches {
		if doc := p.LoadSpecFile(match); doc != nil {
			res = append(res, doc)
		}
	}
	return res
}
