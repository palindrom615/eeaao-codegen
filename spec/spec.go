package spec

import (
	"github.com/palindrom615/eeaao-codegen/plugin"
	"path/filepath"
)

type Spec any

type SpecDir struct {
	Dir    string
	Plugin plugin.Plugin
}

func (s SpecDir) LoadSpecsGlob(glob string) []Spec {
	p := filepath.Join(s.Dir, glob)
	matches, err := filepath.Glob(p)
	if err != nil {
		return nil
	}
	res := []Spec{}
	for _, match := range matches {
		if doc := s.Plugin.LoadSpecFile(match); doc != nil {
			res = append(res, doc)
		}
	}
	return res
}
