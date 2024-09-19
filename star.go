package eeaao_codegen

import (
	json2 "go.starlark.net/lib/json"
	"go.starlark.net/lib/math"
	"go.starlark.net/lib/time"
	"go.starlark.net/repl"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"go.starlark.net/syntax"
	"maps"
	"path/filepath"
)

type starlarkRunner struct {
	thread      *starlark.Thread
	predefined  starlark.StringDict
	fileOptions *syntax.FileOptions
	codeletDir  string
}

func newStarlarkRunner(codeletDir string, eeaaoModule *starlarkstruct.Module) (*starlarkRunner, error) {
	s := &starlarkRunner{
		thread: &starlark.Thread{Name: "main"},
		predefined: starlark.StringDict{
			"json":          json2.Module,
			"math":          math.Module,
			"time":          time.Module,
			"eeaao_codegen": eeaaoModule,
		},
		fileOptions: &syntax.FileOptions{
			Set:               true,
			While:             false,
			TopLevelControl:   false,
			GlobalReassign:    false,
			LoadBindsGlobally: false,
			Recursion:         true,
		},
		codeletDir: codeletDir,
	}
	globals, err := starlark.ExecFileOptions(
		s.fileOptions,
		s.thread,
		filepath.Join(codeletDir, "render.star"),
		nil,
		s.predefined,
	)
	if err != nil {
		return nil, err
	}

	if globals["main"] == nil {
		return nil, syntax.Error{
			Msg: "main function not found in render.star",
		}
	}
	maps.Copy(s.predefined, globals)

	return s, nil
}

// Render runs the main function in the starlark script
func (s *starlarkRunner) Render() (starlark.Value, error) {
	return starlark.Call(
		s.thread,
		s.predefined["main"],
		starlark.Tuple{},
		[]starlark.Tuple{},
	)
}

// RunShell runs the starlark shell
func (s *starlarkRunner) RunShell() {
	repl.REPLOptions(
		s.fileOptions,
		s.thread,
		s.predefined,
	)
}
