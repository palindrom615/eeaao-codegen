package spec

import (
	"github.com/palindrom615/eeaao-codegen/plugin"
	"path/filepath"
)

type SpecDir struct {
	Dir string
}

func (s SpecDir) LoadSpecsGlob(pluginName string, glob string) (res []plugin.SpecData) {
	p := plugin.GetPlugin(pluginName)
	matches, err := filepath.Glob(filepath.Join(s.Dir, glob))
	if err != nil {
		return nil
	}
	for _, match := range matches {
		if doc := p.LoadSpecFile(match); doc != nil {
			res = append(res, doc)
		}
	}
	return res
}
