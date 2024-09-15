// Package codelet declares the exposed functions for go/template and starlark built-ins.
package codelet

import (
	"encoding/json"
	"github.com/Masterminds/sprig"
	json2 "go.starlark.net/lib/json"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"log"
	"maps"
	"text/template"
)

// HelperFuncs is the interface for the helper functions
// that can be used in the template or starlark script
type HelperFuncs interface {
	// LoadSpecsGlob loads specs from a directory with a glob pattern
	// pluginName: the name of the plugin
	// glob: the glob pattern of the spec files from the spec directory
	// returns a map of spec file path and spec content as json encoded string
	LoadSpecsGlob(pluginName string, glob string) (map[string]string, error)
	// RenderFile renders a file with a template
	// filePath: the file path to render
	// templatePath: the template path
	// data: the data to render
	// returns the destination file path
	RenderFile(filePath string, templatePath string, data any) (dst string)
	// WithConfig returns the configuration as a map
	WithConfig() map[string]any
}

func ToTemplateFuncmap(h HelperFuncs) template.FuncMap {
	funcmap := template.FuncMap{
		"loadSpecsGlob": h.LoadSpecsGlob,
		"renderFile":    h.RenderFile,
		"withConfig":    h.WithConfig,
	}
	maps.Copy(funcmap, sprig.FuncMap())
	return funcmap
}

// ToStarlarkModule converts the helper functions to starlarkstruct.Module
// due to the limitation of
func ToStarlarkModule(h HelperFuncs) *starlarkstruct.Module {
	loadSpecsGlob := starlark.NewBuiltin(
		"loadSpecsGlob",
		func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var (
				pluginName, glob starlark.String
			)
			if err := starlark.UnpackArgs("loadSpecsGlob", args, kwargs, "pluginName", &pluginName, "glob", &glob); err != nil {
				return nil, err
			}
			specsLoaded, err := h.LoadSpecsGlob(string(pluginName), string(glob))
			if err != nil {
				return nil, err
			}

			specs := starlark.NewDict(len(specsLoaded))
			for path, spec := range specsLoaded {
				decoded, err := decodeWithStarlarkJson(thread, starlark.String(spec))
				if err != nil {
					log.Printf("Error decoding spec file '%s': %v\n", path, err)
				}
				specs.SetKey(starlark.String(path), decoded)
			}
			return specs, nil
		},
	)
	renderFile := starlark.NewBuiltin(
		"renderFile",
		func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var (
				filePath, templatePath starlark.String
				data                   starlark.Value
			)
			if err := starlark.UnpackArgs("renderFile", args, kwargs, "filePath", &filePath, "templatePath", &templatePath, "data", &data); err != nil {
				return nil, err
			}
			dst := h.RenderFile(string(filePath), string(templatePath), data)
			return starlark.String(dst), nil
		},
	)
	withConfig := starlark.NewBuiltin(
		"withConfig",
		func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			if err := starlark.UnpackArgs("withConfig", args, kwargs); err != nil {
				return nil, err
			}
			conf, err := json.Marshal(h.WithConfig())
			if err != nil {
				return nil, err
			}
			return decodeWithStarlarkJson(thread, starlark.String(conf))
		},
	)
	return &starlarkstruct.Module{
		Name: "eeaao_codegen",
		Members: starlark.StringDict{
			"loadSpecsGlob": loadSpecsGlob,
			"renderFile":    renderFile,
			"withConfig":    withConfig,
		},
	}
}

func decodeWithStarlarkJson(thread *starlark.Thread, value starlark.Value) (starlark.Value, error) {
	decode := json2.Module.Members["decode"].(*starlark.Builtin)
	return starlark.Call(thread, decode, starlark.Tuple{value}, nil)
}
