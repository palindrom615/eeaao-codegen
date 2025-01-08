package eeaao_codegen

import (
	"github.com/palindrom615/eeaao-codegen/starlarkbridge"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"log"
)

// EeaaoStarlarkModule is a starlark module that exposes the App's functions to starlark.
type EeaaoStarlarkModule struct {
	*starlarkstruct.Module
	app *App
}

func NewEeaaoStarlarkModule(app *App) (m *EeaaoStarlarkModule) {
	m = &EeaaoStarlarkModule{
		Module: &starlarkstruct.Module{
			Name:    "eeaao_codegen",
			Members: starlark.StringDict{},
		},
		app: app,
	}
	m.addBuiltinMethod("renderFile", m.renderFile)
	m.addBuiltinMethod("getPlugin", m.getPlugin)
	m.addBuiltinMethod("loadValues", m.loadValues)
	m.addBuiltinMethod("addTemplateFunc", m.addTemplateFunc)
	return
}

// renderFile renders a file using the given template and data.
func (m *EeaaoStarlarkModule) renderFile(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
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
	dst, err := m.app.RenderFile(string(filePath), string(templatePath), d)
	if err != nil {
		return nil, err
	}
	return starlark.String(dst), nil
}

// getPlugin returns the plugin with the given name.
func (m *EeaaoStarlarkModule) getPlugin(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		pluginName starlark.String
	)
	if err := starlark.UnpackArgs("getPlugin", args, kwargs, "pluginName", &pluginName); err != nil {
		return nil, err
	}
	p := m.app.GetPlugin(string(pluginName))
	if p == nil {
		return starlark.None, nil
	}
	return starlarkbridge.NewPluginStarlark(p), nil
}

// loadValues returns the values data from codelet's default values.yaml file and given values file.
func (m *EeaaoStarlarkModule) loadValues(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackArgs("loadValues", args, kwargs); err != nil {
		return nil, err
	}
	return starlarkbridge.ConvertToStarlarkValue(thread, m.app.LoadValues())
}

// addTemplateFunc adds a template function to the template.
func (m *EeaaoStarlarkModule) addTemplateFunc(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		name         starlark.String
		templateFunc starlark.Callable
	)
	if err := starlark.UnpackArgs("addTemplateFunc", args, kwargs, "name", &name, "f", &templateFunc); err != nil {
		return nil, err
	}
	tmplFunc := func(v ...any) (any, error) {
		args := make(starlark.Tuple, len(v))
		kwargs := make([]starlark.Tuple, 0)
		for i, arg := range v {
			var err error
			args[i], err = starlarkbridge.ConvertToStarlarkValue(thread, arg)
			if err != nil {
				return nil, err
			}
		}

		res, err := starlark.Call(thread, templateFunc, args, kwargs)
		if err != nil {
			return nil, err
		}
		return starlarkbridge.ConvertFromStarlarkValue(thread, res)
	}
	m.app.tmpl.AddTemplateFunc(name.GoString(), tmplFunc)
	return starlark.None, nil
}

// addBuiltinMethod adds a builtin method to the module.
func (m *EeaaoStarlarkModule) addBuiltinMethod(fnName string, fn func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)) {
	builtin := starlark.NewBuiltin(
		fnName, fn,
	)
	m.Members[fnName] = builtin
}
