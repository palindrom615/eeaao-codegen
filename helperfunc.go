// Package codelet declares the exposed functions for go/template and starlark built-ins.
package eeaao_codegen

import (
	"encoding/json"
	"github.com/Masterminds/sprig"
	"github.com/palindrom615/eeaao-codegen/plugin"
	"github.com/palindrom615/eeaao-codegen/starlarkbridge"
	json2 "go.starlark.net/lib/json"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"log"
	"maps"
	"text/template"
)

// HelperFuncs defines the helper functions for codelet.
// The functions should be exposed to go/template and starlark built-ins for codelet.
type HelperFuncs interface {
	// LoadSpecFile loads plugin.SpecData from filePath with plugin of pluginName
	// filePath is relative path from spec directory
	LoadSpecFile(pluginName string, filePath string) (plugin.SpecData, error)
	// LoadSpecsGlob loads specs from the given glob pattern in the spec directory.
	// pluginName: the plugin name to load the specs
	// glob: the glob pattern to search for specs
	// returns a map of spec file path and spec content as json encoded string.
	LoadSpecsGlob(pluginName string, glob string) (specs map[string]plugin.SpecData, err error)
	// RenderFile renders a file with the given template and data
	// filePath: the file path to render. The path is relative to the output directory.
	// templatePath: the template path. The path is relative to the ${codeletdir}/templates directory.
	// data: the data to render
	// returns the destination file path.
	RenderFile(filePath string, templatePath string, data any) (dst string, err error)
	// LoadValues returns the values data from codelet's default values.yaml file and given values file.
	LoadValues() (config map[string]any)
	// GetPlugin returns the plugin with the given name.
	GetPlugin(pluginName string) plugin.Plugin
	// Include renders a template with the given data.
	//
	// Drop-in replacement for template pipeline, but with a string return value so that it can be treated as a string in the template.
	//
	// Inspired by [helm include function](https://helm.sh/docs/chart_template_guide/named_templates/#the-include-function)
	Include(templatePath string, data interface{}) (string, error)
}

// ToTemplateFuncmap converts the helper functions into a template.FuncMap
// for use with template.
//
// The resulting FuncMap includes the following functions:
//   - loadSpecFile: HelperFuncs.LoadSpecFile
//   - loadSpecsGlob: HelperFuncs.LoadSpecsGlob
//   - renderFile: HelperFuncs.RenderFile
//   - loadValues: HelperFuncs.LoadValues
//   - include: HelperFuncs.Include
//   - getPlugin: HelperFuncs.GetPlugin
//
// Additionally, it incorporates the sprig.FuncMap for extended functionality.
func ToTemplateFuncmap(h HelperFuncs) template.FuncMap {
	funcmap := template.FuncMap{
		"loadSpecFile":  h.LoadSpecFile,
		"loadSpecsGlob": h.LoadSpecsGlob,
		"renderFile":    h.RenderFile,
		"loadValues":    h.LoadValues,
		"include":       h.Include,
		"getPlugin":     h.GetPlugin,
	}
	maps.Copy(funcmap, sprig.FuncMap())
	return funcmap
}

// ToStarlarkModule exposes the helper functions to starlarkstruct.Module.
//
// The module provides the following functions:
//   - loadSpecFile(pluginName: str, filepath: str) -> any
//   - loadSpecsGlob(pluginName: str, glob: str) -> dict[str, any]
//   - renderFile(filePath: str, templatePath: str, data: any) -> str
//   - loadValues() -> dict[str, any]
//   - getPlugin(pluginName: str) -> EeaaoPlugin
//
// due to the limitation of interoperability between Go and Starlark, JSON encoding is used internally.
//
// For example, HelperFuncs.LoadValues() returns a map[string]any. When
// `eeaao_codegen.loadValues` is called in starlark script, the map[string]any is first encoded
// using json.Marshal and then decoded as starlark.Dict by `json.decode`.
func ToStarlarkModule(h HelperFuncs) *starlarkstruct.Module {
	loadSpecFile := starlark.NewBuiltin(
		"loadSpecFile",
		func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var (
				pluginName, filepath starlark.String
			)
			if err := starlark.UnpackArgs("loadSpecFile", args, kwargs, "pluginName", &pluginName, "filepath", &filepath); err != nil {
				return nil, err
			}
			spec, err := h.LoadSpecFile(string(pluginName), string(filepath))
			if err != nil {
				return nil, err
			}
			return convertToStarlarkValue(thread, spec)
		},
	)
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
				decoded, err := convertToStarlarkValue(thread, spec)
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
			d, err := convertFromStarlarkValue(thread, data)
			if err != nil {
				log.Printf("Error decoding starlark injected data: %v\n%v\n", data, err)
				return nil, err
			}
			dst, err := h.RenderFile(string(filePath), string(templatePath), d)
			if err != nil {
				return nil, err
			}
			return starlark.String(dst), nil
		},
	)
	loadValues := starlark.NewBuiltin(
		"loadValues",
		func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			if err := starlark.UnpackArgs("loadValues", args, kwargs); err != nil {
				return nil, err
			}
			return convertToStarlarkValue(thread, h.LoadValues())
		},
	)
	getPlugin := starlark.NewBuiltin(
		"getPlugin",
		func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var (
				pluginName starlark.String
			)
			if err := starlark.UnpackArgs("getPlugin", args, kwargs, "pluginName", &pluginName); err != nil {
				return nil, err
			}
			p := h.GetPlugin(string(pluginName))
			if p == nil {
				return starlark.None, nil
			}
			return starlarkbridge.NewPluginStarlark(p), nil
		},
	)
	return &starlarkstruct.Module{
		Name: "eeaao_codegen",
		Members: starlark.StringDict{
			"loadSpecFile":  loadSpecFile,
			"loadSpecsGlob": loadSpecsGlob,
			"renderFile":    renderFile,
			"loadValues":    loadValues,
			"getPlugin":     getPlugin,
		},
	}
}

func convertToStarlarkValue(thread *starlark.Thread, value any) (starlark.Value, error) {
	specStr, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return decodeWithStarlarkJson(thread, starlark.String(specStr))
}

func convertFromStarlarkValue(thread *starlark.Thread, value starlark.Value) (map[string]any, error) {
	encoded, err := encodeWithStarlarkJson(thread, value)
	if err != nil {
		return nil, err
	}
	d := make(map[string]any)
	err = json.Unmarshal([]byte(encoded.(starlark.String)), &d)
	if err != nil {
		log.Printf("Error decoding starlark injected data: %v\n%v\n", encoded, err)
		return nil, err
	}
	return d, nil
}

func decodeWithStarlarkJson(thread *starlark.Thread, value starlark.Value) (starlark.Value, error) {
	decode := json2.Module.Members["decode"].(*starlark.Builtin)
	return starlark.Call(thread, decode, starlark.Tuple{value}, nil)
}

func encodeWithStarlarkJson(thread *starlark.Thread, value starlark.Value) (starlark.Value, error) {
	encode := json2.Module.Members["encode"].(*starlark.Builtin)
	return starlark.Call(thread, encode, starlark.Tuple{value}, nil)
}
