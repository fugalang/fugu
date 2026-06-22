package recursive

import (
	"github.com/fugalang/fugu/internal/ast"
	"github.com/fugalang/fugu/internal/token"
)

func (p *Parser) expr() int {
	left := p.term()

	for p.match(token.ADD, token.SUB) {
		op := p.pastTk

		right := p.term()

		return p.addNode(ast.Node{
			Type:  ast.Binary,
			Data1: left,
			Data2: right,
			Data3: int(ast.Op(op.Kind)),
		})
	}
	return left
}

func (p *Parser) term() int {
	left := p.factor()

	for p.match(token.MUL, token.DIV, token.MOD) {
		op := p.pastTk

		right := p.factor()

		return p.addNode(ast.Node{
			Type:  ast.Binary,
			Data1: left,
			Data2: right,
			Data3: int(ast.Op(op.Kind)),
		})
	}

	return left
}

func (p *Parser) factor() int {
	if p.match(token.SUB) {
		op := p.pastTk
		expr := p.factor()

		return p.addNode(ast.Node{
			Type:  ast.Unary,
			Data1: expr,
			Data3: int(ast.Op(op.Kind)),
		})
	}

	if p.match(token.G_NUMBER) {

		return p.addNode(ast.Node{
			Type:  ast.Literal,
			Data1: p.addString(p.VauleToken()),
		})
	}

	if p.match(token.L_PAREN) {
		expr := p.expr()
		p.expect(token.R_PAREN)
		return expr
	}

	return -1
}

func (p *Parser) parsModule() {
	if p.eat(token.MODULE) {
		if p.eat(token.IDENTIFIER) {
			nameModule := p.VauleToken()
			i := p.addString(nameModule)
			if nameModule == "main" {
				p.addNode(ast.Node{
					Type:  ast.Module,
					Data1: 1,
					Data2: i,
				})
			} else {
				p.addNode(ast.Node{
					Type:  ast.Module,
					Data1: 2,
					Data2: i,
				})
				return
			}
		} else {
			return
		}
	}
}
