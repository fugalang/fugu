package ast

import (
	"github.com/fugalang/fugu/internal/types"

	"github.com/fugalang/fugu/internal/ast"
	"github.com/fugalang/fugu/internal/parser"
)

type AstContext struct {
	Scopes []map[string]types.Type
}

func Analysis(a *ast.AstArena, pars *parser.Parser) {
	Walk(0, a)
}

func Walk(i int, a *ast.AstArena) {
	n := a.Nodes[i]

	switch n.Type {
	case ast.BinaryExpr:
		left := n.Data1
		right := n.Data2

		Walk(left, a)
		Walk(right, a)

	case ast.UnaryExpr:
		Walk(n.Data1, a)

	case ast.Literal:

	}
}
