package starlarkbridge

import (
	"github.com/palindrom615/eeaao-codegen/plugin"
	"go.starlark.net/starlark"
	"hash/fnv"
)

// NewPluginStarlark returns a starlark.HasAttrs for the given plugin.Plugin
//
// in starlark, you can access below methods from the value returned by this function
//
//	spec_from_file = plugin.loadSpecFile("path/to/spec.yaml")
//	spec_from_url = plugin.loadSpecUrl("https://example.com/spec.yaml")
func NewPluginStarlark(plugin plugin.Plugin) starlark.HasAttrs {
	loadSpecFile := &pluginMethodStarlark{
		plugin: plugin,
		name:   "loadSpecFile",
		callInternal: func(thread *starlark.Thread, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var (
				path starlark.String
			)
			if err := starlark.UnpackArgs("loadSpecFile", args, kwargs, "path", &path); err != nil {
				return nil, err
			}
			spec, err := plugin.LoadSpecFile(string(path))
			if err != nil {
				return nil, err
			}
			return ConvertToStarlarkValue(thread, spec)
		},
	}
	loadSpecUrl := &pluginMethodStarlark{
		plugin: plugin,
		name:   "loadSpecUrl",
		callInternal: func(thread *starlark.Thread, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var (
				url starlark.String
			)
			if err := starlark.UnpackArgs("loadSpecUrl", args, kwargs, "url", &url); err != nil {
				return nil, err
			}
			spec, err := plugin.LoadSpecUrl(string(url))
			if err != nil {
				return nil, err
			}
			return ConvertToStarlarkValue(thread, spec)
		},
	}
	return &pluginStarlark{
		plugin:       plugin,
		loadSpecFile: loadSpecFile,
		loadSpecUrl:  loadSpecUrl,
	}
}

// pluginStarlark is the starlark interface for plugin.Plugin
// It implements starlark.Value and starlark.HasAttrs interfaces.
// It provides the following methods:
//   - loadSpecFile: loadSpecFile(url: string) -> any
//   - loadSpecUrl: loadSpecUrl(url: string) -> any
type pluginStarlark struct {
	plugin       plugin.Plugin
	loadSpecFile starlark.Callable
	loadSpecUrl  starlark.Callable
}

func (p *pluginStarlark) String() string {
	return "EEAAO " + p.plugin.Name() + " Plugin"
}

func (p *pluginStarlark) Type() string {
	return "EeaaoPlugin"
}

func (p *pluginStarlark) Freeze() {

}

func (p *pluginStarlark) Truth() starlark.Bool {
	return true
}

func (p *pluginStarlark) Hash() (uint32, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(p.plugin.Name()))
	if err != nil {
		return 0, err
	}
	return h.Sum32(), nil
}

func (p *pluginStarlark) Attr(name string) (starlark.Value, error) {
	switch name {
	case "loadSpecFile":
		return p.loadSpecFile, nil
	case "loadSpecUrl":
		return p.loadSpecUrl, nil
	}
	return nil, nil
}

func (p *pluginStarlark) AttrNames() []string {
	return []string{"loadSpecFile", "loadSpecUrl"}
}

// pluginMethodStarlark is the starlark interface for plugin.Plugin methods
type pluginMethodStarlark struct {
	plugin       plugin.Plugin
	name         string
	callInternal func(thread *starlark.Thread, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)
}

func (p *pluginMethodStarlark) String() string {
	return "EEAAO " + p.plugin.Name() + " Plugin Method " + p.name
}

func (p *pluginMethodStarlark) Type() string {
	return "EeaaoPluginMethod"
}

func (p *pluginMethodStarlark) Freeze() {

}

func (p *pluginMethodStarlark) Truth() starlark.Bool {
	return true
}

func (p *pluginMethodStarlark) Hash() (uint32, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(p.plugin.Name()))
	if err != nil {
		return 0, err
	}
	_, err = h.Write([]byte(p.name))
	if err != nil {
		return 0, err
	}
	return h.Sum32(), nil
}

func (p *pluginMethodStarlark) Name() string {
	return p.name
}

func (p *pluginMethodStarlark) CallInternal(thread *starlark.Thread, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return p.callInternal(thread, args, kwargs)
}
