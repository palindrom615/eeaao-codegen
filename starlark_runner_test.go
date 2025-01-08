package eeaao_codegen

import (
	"go.starlark.net/starlarkstruct"
	"testing"
)

func TestStarlarkRunner_newStarlarkRunner(t *testing.T) {
	codeletDir := "test/codelet"
	eeaaoModule := &starlarkstruct.Module{}
	_, err := newStarlarkRunner(codeletDir, eeaaoModule)
	if err != nil {
		t.Errorf("Error creating starlark runner: %v", err)
	}
}
