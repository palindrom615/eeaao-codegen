package eeaao_codegen

import (
	json2 "go.starlark.net/lib/json"
	"go.starlark.net/lib/math"
	"go.starlark.net/lib/time"
	"go.starlark.net/repl"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"go.starlark.net/syntax"
)

type starlarkRunner struct {
	thread      *starlark.Thread
	predefined  starlark.StringDict
	fileOptions *syntax.FileOptions
}

func newStarlarkRunner() *starlarkRunner {
	return &starlarkRunner{
		thread: &starlark.Thread{Name: "main"},
		predefined: starlark.StringDict{
			"json": json2.Module,
			"math": math.Module,
			"time": time.Module,
		},
		fileOptions: &syntax.FileOptions{
			Set:               true,
			While:             true,
			TopLevelControl:   true,
			GlobalReassign:    true,
			LoadBindsGlobally: false,
			Recursion:         true,
		},
	}
}

func (s *starlarkRunner) addModule(name string, module *starlarkstruct.Module) {
	s.predefined[name] = module
}

func (s *starlarkRunner) runFile(file string) (starlark.StringDict, error) {
	return starlark.ExecFileOptions(
		s.fileOptions,
		s.thread,
		file,
		nil,
		s.predefined,
	)
}

func (s *starlarkRunner) runShell() {
	repl.REPLOptions(
		s.fileOptions,
		s.thread,
		s.predefined,
	)
}
