package eeaao_codegen

import (
	"github.com/Masterminds/sprig"
	"github.com/palindrom615/eeaao-codegen/starlarkbridge"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"log"
	"maps"
	"text/template"
)

// ToTemplateFuncmap converts the helper functions into a template.FuncMap
// for use with template.
//
// The resulting FuncMap includes the following functions:
//   - renderFile: App.RenderFile
//   - loadValues: App.LoadValues
//   - include: App.Include
//   - getPlugin: App.GetPlugin
//
// Additionally, it incorporates the sprig.FuncMap for extended functionality.
func ToTemplateFuncmap(a *App) template.FuncMap {
	funcmap := template.FuncMap{
		"renderFile": a.RenderFile,
		"loadValues": a.LoadValues,
		"include":    a.Include,
		"getPlugin":  a.GetPlugin,
	}
	maps.Copy(funcmap, sprig.FuncMap())
	return funcmap
}

// ToStarlarkModule exposes the App's functions to starlarkstruct.Module.
//
// The module provides the following functions:
//   - renderFile(filePath: str, templatePath: str, data: any) -> str
//   - loadValues() -> dict[str, any]
//   - getPlugin(pluginName: str) -> EeaaoPlugin
//
// due to the limitation of interoperability between Go and Starlark, JSON encoding is used internally.
//
// For example, App.LoadValues() returns a map[string]any. When
// `eeaao_codegen.loadValues` is called in starlark script, the map[string]any is first encoded
// using json.Marshal and then decoded as starlark.Dict by `json.decode`.
func ToStarlarkModule(app *App) *starlarkstruct.Module {
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
			d, err := starlarkbridge.ConvertFromStarlarkValue(thread, data)
			if err != nil {
				log.Printf("Error decoding starlark injected data: %v\n%v\n", data, err)
				return nil, err
			}
			dst, err := app.RenderFile(string(filePath), string(templatePath), d)
			if err != nil {
				return nil, err
			}
			return starlark.String(dst), nil
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
			p := app.GetPlugin(string(pluginName))
			if p == nil {
				return starlark.None, nil
			}
			return starlarkbridge.NewPluginStarlark(p), nil
		},
	)
	loadValues := starlark.NewBuiltin(
		"loadValues",
		func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			if err := starlark.UnpackArgs("loadValues", args, kwargs); err != nil {
				return nil, err
			}
			return starlarkbridge.ConvertToStarlarkValue(thread, app.LoadValues())
		},
	)
	return &starlarkstruct.Module{
		Name: "eeaao_codegen",
		Members: starlark.StringDict{
			"renderFile": renderFile,
			"loadValues": loadValues,
			"getPlugin":  getPlugin,
		},
	}
}
