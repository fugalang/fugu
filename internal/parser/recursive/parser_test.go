package recursive

import (
	"testing"

	"github.com/k0kubun/pp/v3"
)

func TestParserAst(t *testing.T) {
	pars := New([]byte("module main"), "main.fg")

	ast := pars.Parse()
	pp.Println(ast.Nodes[0])
}
