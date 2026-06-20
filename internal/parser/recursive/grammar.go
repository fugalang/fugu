package recursive

import (
	"github.com/fugalang/fugu/internal/ast"
	"github.com/fugalang/fugu/internal/token"
)

func (p *Parser) parsModule() {
	if p.eat(token.MODULE) {
		if p.eat(token.IDENTIFIER) {
			i := len(p.Ast.Strings)
			p.Ast.Strings = append(p.Ast.Strings, string(p.curTk.Literal(&p.input)))

			p.Ast.Nodes = append(p.Ast.Nodes, ast.Node{
				Type: ast.Module,

				Data1: i,
			})
		} else {
			return
		}
	}
}
