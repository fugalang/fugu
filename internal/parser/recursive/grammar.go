package recursive

import (
	"strconv"

	"github.com/fugalang/fugu/internal/ast"
	"github.com/fugalang/fugu/internal/token"
)

func (p *Parser) expr() int {
	left := p.term()

	for p.match(token.ADD, token.SUB) {
		op := p.pastTk
		right := p.term()

		left = p.AddNode(ast.Node{
			Type:  ast.Binary,
			Data1: left,
			Data2: right,
			Data3: int(ast.Op(op.Kind)),
		})
	}

	return left
}

func (p *Parser) term() int {
	left := p.power()

	for p.match(token.MUL, token.DIV, token.MOD) {
		op := p.pastTk

		right := p.power()

		left = p.AddNode(ast.Node{
			Type:  ast.Binary,
			Data1: left,
			Data2: right,
			Data3: int(ast.Op(op.Kind)),
		})
	}

	return left
}

func (p *Parser) power() int {
	left := p.factor()

	if p.match(token.POW) {
		op := p.pastTk
		right := p.power()

		return p.AddNode(ast.Node{
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

		return p.AddNode(ast.Node{
			Type:  ast.Unary,
			Data1: expr,
			Data3: int(ast.Op(op.Kind)),
		})
	}

	if p.match(token.G_NUMBER) {
		switch p.fkTk {
		case token.INTEGER:
			v, _ := strconv.ParseInt(p.VauleToken(), 10, 64)
			return p.AddNode(ast.Node{
				Type: ast.Literal,
				Data1: p.AddString(ast.Value{
					Type: ast.Int,
					I64:  v,
				}),
			})

		case token.IMAGINARY:
			v, _ := strconv.ParseComplex(p.VauleToken(), 128)
			return p.AddNode(ast.Node{
				Type: ast.Literal,
				Data1: p.AddString(ast.Value{
					Type: ast.Complex,
					C128: v,
				}),
			})

		case token.FLOATING:
			v, _ := strconv.ParseFloat(p.VauleToken(), 64)
			return p.AddNode(ast.Node{
				Type: ast.Literal,
				Data1: p.AddString(ast.Value{
					Type: ast.Float,
					F64:  v,
				}),
			})
		}
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
			i := p.AddString(ast.Value{
				Type: ast.String,
				STR:  nameModule,
			})
			if nameModule == "main" {
				p.AddNode(ast.Node{
					Type:  ast.Module,
					Data1: 1,
					Data2: i,
				})
			} else {
				p.AddNode(ast.Node{
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

func (p *Parser) parsLet() {
	if p.match(token.LET) {
		if p.match(token.IDENTIFIER) {
			name := p.VauleToken()
			if p.match(token.COLON) {
				if p.match(token.MUT) {
					if p.match(token.IDENTIFIER) {
						if p.match(token.ASSIGN) {
							if p.match(token.G_LITERAL) {
								p.AddNode(ast.Node{
									Type: ast.Let,

									Data1: p.AddString(ast.Value{
										Type: ast.String,
										STR:  name,
									}),
									Data2: 1,
									Data3: 1,
								})
							} else {

							}
						} else {

						}
					} else {

					}
				} else {
					if p.match(token.CONST) {

					} else {

					}
				}
			} else {

			}
		} else {

		}
	} else {

	}
}
