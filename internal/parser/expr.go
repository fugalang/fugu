package parser

import (
	"github.com/fugalang/fugu/internal/ast"
	. "github.com/fugalang/fugu/internal/token"
	"github.com/fugalang/fugu/pkg/helper"
)

const (
	Lowest   = iota
	Or       // ||
	And      // &&
	Equality // == !=
	Compare  // < > <= >=
	Sum      // + -
	Product  // * / %
	Power    // ^
	Prefix   // -x
	Call     // f()
)

func precedence(k Kind) int {
	switch k {
	case OR:
		return Or
	case AND:
		return And

	case EQ, NEQ:
		return Equality

	case LT, GT, LE, GE:
		return Compare

	case ADD, SUB:
		return Sum

	case MUL, DIV, MOD:
		return Product

	case POW:
		return Power

	default:
		return Lowest
	}
}

func (p *Parser) parseExpr(pre int) int {
	left := p.parsePrefix()

	for {
		op := p.curTk.Kind
		if op == EOF {
			break
		}

		pred := precedence(op)
		if pred <= pre {
			break
		}

		p.next()

		right := p.parseExpr(pred)

		left = p.addNode(ast.Node{
			Type:  ast.BinaryExpr,
			Data1: int(ast.Op(op)),
			Data2: left,
			Data3: right,
		})
	}

	return left
}

func (p *Parser) parsePrefix() int {
	switch p.curTk.Kind {

	case SUB:
		p.next()
		expr := p.parseExpr(Prefix)

		return p.addNode(ast.Node{
			Type:  ast.UnaryExpr,
			Data1: int(ast.OpNeg),
			Data2: expr,
		})

	case INTEGER, FLOATING, CHARACTER, IMAGINARY:
		return p.parseLiteral()

	case IDENTIFIER:
		return p.parseIdent()

	case L_PAREN:
		p.next()
		expr := p.parseExpr(Lowest)
		p.match(R_PAREN)
		return expr

	case STRING:
		return p.parseTemplate()

	default:
		return p.addNode(ast.Node{Type: ast.Invalid})
	}
}

func (p *Parser) parseLiteral() int {
	tk := p.curTk
	p.next()

	switch tk.Kind {

	case INTEGER:
		i := p.addValue(ast.Value{
			Type: ast.Int,
			I64:  helper.StringIn64(string(tk.Literal(&p.lex.Input))),
		})

		return p.addNode(ast.Node{
			Type:  ast.Literal,
			Data1: i,
		})

	case FLOATING:
		i := p.addValue(ast.Value{
			Type: ast.Float,
			F64:  helper.StringFloat64(string(tk.Literal(&p.lex.Input))),
		})

		return p.addNode(ast.Node{
			Type:  ast.Literal,
			Data1: i,
		})

	case IMAGINARY:
		i := p.addValue(ast.Value{
			Type: ast.Complex,
			C128: helper.StringComplex128(string(tk.Literal(&p.lex.Input))),
		})

		return p.addNode(ast.Node{
			Type:  ast.Literal,
			Data1: i,
		})

	case CHARACTER:
		r := []rune(string(tk.Literal(&p.lex.Input)))

		i := p.addValue(ast.Value{
			Type: ast.Char,
			I64:  int64(r[0]),
		})

		return p.addNode(ast.Node{
			Type:  ast.Literal,
			Data1: i,
		})

	case RAW_STRING:
		i := p.addValue(ast.Value{
			Type: ast.String,
			S8:   string(tk.Literal(&p.lex.Input)),
		})

		return p.addNode(ast.Node{
			Type:  ast.Literal,
			Data1: i,
		})
	}

	return p.addNode(ast.Node{
		Type: ast.Invalid,
	})
}

func (p *Parser) parseIdent() int {
	tk := p.curTk
	p.next()

	idx := p.addValue(ast.Value{
		Type: ast.String,
		S8:   string(tk.Literal(&p.lex.Input)),
	})

	return p.addNode(ast.Node{
		Type:  ast.Ident,
		Data1: idx,
	})
}

func (p *Parser) parseTemplate() int {
	node := ast.Node{
		Type: ast.Template,
	}

	for {
		if p.curTk.Kind == STRING {
			p.addNode(ast.Node{
				Type: ast.Template,
				Data1: p.addValue(ast.Value{
					Type: ast.String,
					S8:   p.VauleToken(),
				}),
			})

			p.next()
			continue
		}

		if p.curTk.Kind == L_BRACE {
			p.next()

			expr := p.parseExpr(Lowest)

			p.match(R_BRACE)

			p.addNode(ast.Node{
				Type:  ast.Template,
				Data2: expr,
			})

			continue
		}

		break
	}

	return p.addNode(node)
}
