package recursive

import (
	"fmt"

	"github.com/fugalang/fugu/internal/ast"
	"github.com/fugalang/fugu/internal/diagnostics"
	"github.com/fugalang/fugu/internal/diagnostics/errors"
	"github.com/fugalang/fugu/internal/lexer"
	"github.com/fugalang/fugu/internal/token"
)

type Parser struct {
	lex *lexer.Lexer

	Tokens []token.Token
	curTk  token.Token
	pastTk token.Token
	fkTk   token.Kind
	pos    int

	pn int
	si int

	Ast *ast.AstArena

	da *diagnostics.Arena
}

func New(input []byte, fileName string) *Parser {
	p := &Parser{
		Tokens: make([]token.Token, 1024),
		da: &diagnostics.Arena{
			Source: string(input),
		},
		Ast: &ast.AstArena{
			Nodes: make([]ast.Node, 1024),
			Value: make([]ast.Value, 512),
		},
	}
	p.lex = lexer.New(input, fileName, p.da)
	p.next()
	p.pos = 0
	return p
}

func (p *Parser) Parse() *ast.AstArena {
	for p.curTk.Kind != token.EOF {
		switch p.curTk.Kind {
		// decl
		case token.MODULE:
			p.parsModule()
		// delc + stmt
		case token.LET:
			p.parsLet()
		}
	}
	return p.Ast
}

func (p *Parser) next() *Parser {
	for {
		tk := p.lex.NextToken()

		if tk.Kind == token.SPACING ||
			tk.Kind == token.COMMENT ||
			tk.Kind == token.M_COMMENT {
			continue
		}

		if tk.Kind == token.EOF && p.curTk.Kind == token.EOF {
			p.pastTk = p.curTk
			return p
		}

		p.pastTk = p.curTk
		p.curTk = tk
		p.fkTk = tk.Kind
		p.curTk.Kind = token.Group(p.curTk.Kind)

		p.Tokens[p.pos] = tk
		p.pos++
		return p
	}
}

func (p *Parser) AddString(v ast.Value) int {
	p.Ast.Value[p.si] = v
	p.si++
	return p.si - 1
}

func (p *Parser) AddNode(n ast.Node) int {
	p.Ast.Nodes[p.pn] = n
	p.pn++
	return p.pn - 1
}

func (p *Parser) match(kinds ...token.Kind) bool {
	for _, kind := range kinds {
		if p.curTk.Kind == kind {
			p.next()
			return true
		}
	}
	return false
}

func (p *Parser) expect(kind token.Kind) bool {
	if p.curTk.Kind == kind {
		p.next()
		return true
	} else {
		p.da.Add(errors.Errors[5].IU("PARSER", []string{
			fmt.Sprintf("ожидалось: %s было получино: %s", kind.String(), p.curTk.Kind.String()),
		}))
		return false
	}
}

func (p *Parser) eat(kind token.Kind) bool {
	if p.curTk.Kind != kind {
		p.da.Add(errors.Errors[5].IU("PARSER", []string{
			fmt.Sprintf("ожидалось: %s было получино: %s", kind.String(), p.curTk.Kind.String()),
		}))
		return false
	}
	p.next()
	return true
}

func (p *Parser) VauleToken() string {
	return string(p.Tokens[p.pos].Literal(&p.lex.Input))
}
