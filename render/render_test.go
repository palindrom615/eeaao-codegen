package render_test

import (
	"github.com/palindrom615/eeaao-codegen/render"
	"testing"
)

func TestRender(t *testing.T) {
	render.Render("../example/spec", "../example/codelet", "../example/build")
}
