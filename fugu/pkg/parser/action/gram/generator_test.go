package gram

import (
	"strings"
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

func TestRenderMaps(t *testing.T) {
	src := `
%start Expr
%token GLITERAL EOF

Expr:
	GLITERAL => KindLiteral
`

	g, err := Parse(src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	table, err := Build(g)
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	out, err := RenderMaps(table, RenderOptions{
		PackageName: "action",
		Imports: []string{
			`"fugu/pkg/parser/action/ast"`,
			`. "fugu/pkg/token"`,
		},
	})
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	for _, want := range []string{
		"var ActionSrc = &map[int]map[TokenKind]ActionStruct",
		"GLITERAL: Sh(",
		"Red(1, ast.KindLiteral)",
		"EOF: Acc()",
		"var GotoSrc = &map[int]map[ast.NodeKind]int",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("rendered maps do not contain %q:\n%s", want, out)
		}
	}
}
