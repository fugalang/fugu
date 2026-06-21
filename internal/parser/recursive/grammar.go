package recursive

import (
	"github.com/fugalang/fugu/internal/ast"
	"github.com/fugalang/fugu/internal/token"
)

func (p *Parser) parsModule() {
	if p.eat(token.MODULE) {
		if p.eat(token.IDENTIFIER) {
			nameModule := string(p.Tokens[p.pos].Literal(&p.input))
			p.addString(nameModule)
			if nameModule == "main" {
				p.addNode(ast.Node{
					Type:  ast.Module,
					Data1: 1,
					Data2: p.si - 1,
				})
			} else {
				p.addNode(ast.Node{
					Type:  ast.Module,
					Data1: 2,
					Data2: p.si - 1,
				})
				return
			}
		} else {
			return
		}
	}
}
