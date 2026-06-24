package gram

import (
	"testing"
)

func TestBuildExpressionGrammar(t *testing.T) {
	src := `
%start Expr
%token GLITERAL PLUS STAR EOF
%left PLUS
%left STAR

Expr:
	Expr PLUS Expr => Expr
	| Expr STAR Expr => Expr
	| GLITERAL => Literal
`

	g, err := Parse(src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	table, err := Build(g)
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	hasAccept := false
	for _, row := range table.Actions {
		for tk, act := range row {
			if tk == EOF && act.Kind == ActionAccept {
				hasAccept = true
			}
		}
	}
	if !hasAccept {
		t.Fatal("table has no EOF accept")
	}
}
