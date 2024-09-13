package eeaao_codegen_test

import (
	"github.com/palindrom615/eeaao-codegen"
	"testing"
)

func TestRender(t *testing.T) {
	c := &eeaao_codegen.App{
		SpecDir:    "./example/spec",
		OutDir:     "./example/build",
		CodeletDir: "./example/codelet",
	}
	eeaao_codegen.Render(c)
}
