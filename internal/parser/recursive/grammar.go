package recursive

import (
	"strconv"

	"github.com/fugalang/fugu/internal/ast"
	"github.com/fugalang/fugu/internal/token"
	"github.com/fugalang/fugu/pkg/helper"
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

	if p.match(token.L_PAREN) {
		expr := p.expr()
		p.expect(token.R_PAREN)
		return expr
	}

	if !p.match(token.G_NUMBER, token.G_STRING, token.G_LITERAL, token.G_ARITHMETIC) {
		return -1
	}

	switch p.fkTk {
	case token.INTEGER:
		v, _ := strconv.ParseInt(p.VauleToken(), 10, 64)
		return p.makeInt(v)

	case token.FLOATING:
		v, _ := strconv.ParseFloat(p.VauleToken(), 64)
		return p.makeFloat(v)

	case token.IMAGINARY:
		v, _ := strconv.ParseComplex(p.VauleToken(), 128)
		return p.makeComplex(v)

	case token.STRING:
		return p.makeString(p.VauleToken())

	case token.RAW_STRING:
		return p.makeString(p.VauleToken())

	case token.T_STRING:
		return p.makeTemplateString(p.VauleToken())

	case token.CHARACTER:
		return p.makeChar(p.VauleToken())

	case token.IDENTIFIER:
		return p.makeIdent(p.VauleToken())

	default:
		return -1
	}
}

func (p *Parser) makeInt(v int64) int {
	return p.AddNode(ast.Node{
		Type: ast.Literal,
		Data1: p.AddString(ast.Value{
			Type: ast.Int,
			I64:  v,
		}),
	})
}

func (p *Parser) makeFloat(v float64) int {
	return p.AddNode(ast.Node{
		Type: ast.Literal,
		Data1: p.AddString(ast.Value{
			Type: ast.Float,
			F64:  v,
		}),
	})
}

func (p *Parser) makeComplex(v complex128) int {
	return p.AddNode(ast.Node{
		Type: ast.Literal,
		Data1: p.AddString(ast.Value{
			Type: ast.Complex,
			C128: v,
		}),
	})
}

func (p *Parser) makeTemplateString(v string) int {
	return p.makeString(v)
}

func (p *Parser) makeString(v string) int {
	return p.AddNode(ast.Node{
		Type: ast.Literal,
		Data1: p.AddString(ast.Value{
			Type: ast.String,
			S8:   v,
		}),
	})
}

func (p *Parser) makeChar(v string) int {
	return p.AddNode(ast.Node{
		Type: ast.Literal,
		Data1: p.AddString(ast.Value{
			Type: ast.Char,
			S8:   v,
		}),
	})
}

func (p *Parser) makeIdent(v string) int {
	return p.AddNode(ast.Node{
		Type: ast.Ident,
		Data1: p.AddString(ast.Value{
			Type: ast.String,
			S8:   v,
		}),
	})
}

func (p *Parser) parseTemplateString(raw string) int {
	parts := []ast.TemplatePart{}

	i := 0
	start := 0

	for i < len(raw) {

		// found ${
		if i+1 < len(raw) && raw[i] == '$' && raw[i+1] == '{' {

			// push string before ${
			if start < i {
				parts = append(parts, ast.TemplatePart{
					IsExpr: false,
					Str:    raw[start:i],
				})
			}

			i += 2 // skip ${
			start = i

			// parse expression until }
			exprStr, newPos := p.extractUntilBrace(raw, i)

			// временно: подставим expr parser через lexer
			exprNode := p.parseInlineExpr(exprStr)

			parts = append(parts, ast.TemplatePart{
				IsExpr: true,
				Expr:   exprNode,
			})

			i = newPos
			start = i
			continue
		}

		i++
	}

	// tail string
	if start < len(raw) {
		parts = append(parts, ast.TemplatePart{
			IsExpr: false,
			Str:    raw[start:],
		})
	}

	return p.AddNode(ast.Node{
		Type:  ast.Template,
		Data1: p.storeTemplate(parts),
	})
}

func (p *Parser) parseInlineExpr(src string) int {
	lexer := NewLexer(src)
	subParser := NewParser(lexer)

	return subParser.expr()
}

func (p *Parser) extractUntilBrace(s string, i int) (string, int) {
	start := i
	depth := 1

	for i < len(s) {
		if s[i] == '{' {
			depth++
		}
		if s[i] == '}' {
			depth--
			if depth == 0 {
				return s[start:i], i + 1
			}
		}
		i++
	}

	return s[start:], i
}

func (p *Parser) parsModule() {
	if p.eat(token.MODULE) {
		if p.eat(token.IDENTIFIER) {
			nameModule := p.VauleToken()
			i := p.AddString(ast.Value{
				Type: ast.String,
				S8:   nameModule,
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
	if !p.eat(token.LET) {
		return
	}

	if !p.eat(token.IDENTIFIER) {
		return
	}

	name := p.VauleToken()

	if !p.eat(token.COLON) {
		return
	}

	isMut := p.eat(token.MUT)

	if !p.eat(token.IDENTIFIER) {
		return
	}

	if !p.eat(token.ASSIGN) {
		return
	}

	value := p.expr() // 💥 ВСЕГДА expr()

	p.AddNode(ast.Node{
		Type: ast.Let,
		Data1: p.AddString(ast.Value{
			Type: ast.String,
			S8:   name,
		}),
		Data2: value,
		Data3: helper.BoolInt(isMut),
	})
}
