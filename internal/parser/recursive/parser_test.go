package recursive

import (
	"testing"

	"github.com/fugalang/fugu/internal/ast"
)

func TestParserAst(t *testing.T) {
	pars := New([]byte("module main"), "main.fg")

	a := pars.Parse()
	as := ast.Node{
		Type:  ast.Module,
		Data2: 0,
		Data1: 1,
	}

	if !(a.Nodes[0] == as) {
		t.Fatalf("Ошибка парсинга module")
	}
}
