package eeaao_codegen_test

import (
	"github.com/palindrom615/eeaao-codegen"
	"testing"
)

func TestRender(t *testing.T) {
	a := eeaao_codegen.NewApp(
		"./example/spec",
		"./example/build/javascript",
		"./example/codelets/javascript",
		"./example/config.yaml",
	)
	a.Render()
}
