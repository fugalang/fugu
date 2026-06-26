package parser

import (
	"github.com/fugalang/fugu/internal/ast"
	. "github.com/fugalang/fugu/internal/token"
)

func (p *Parser) module() int {
	p.next()
	if p.match(IDENTIFIER) {
		return p.addNode(
			ast.Node{
				Type: ast.ModuleDecl,
				Data1: p.addValue(
					ast.Value{
						Type: ast.String,
						S8:   p.VaulePastToken(),
					},
				),
			},
		)
	}
	return -1 // TODO сделать ошибку
}
