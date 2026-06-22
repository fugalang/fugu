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
	input []byte
	lex   *lexer.Lexer

	Tokens []token.Token
	curTk  token.Token
	pastTk token.Token
	pos    int

	pn int
	si int

	Ast *ast.AstArena

	da *diagnostics.DiagnosticArena
}

func New(input []byte, fileName string) *Parser {
	p := &Parser{
		input:  input,
		Tokens: make([]token.Token, 1024),
		da: &diagnostics.DiagnosticArena{
			Source: string(input),
		},
		Ast: &ast.AstArena{
			Nodes:   make([]ast.Node, 1024),
			Strings: make([]string, 512),
		},
	}
	p.lex = lexer.New(input, fileName, p.da)
	p.advance()
	p.pos = 0
	return p
}

func (p *Parser) Parse() *ast.AstArena {
	for p.curTk.Kind != token.EOF {
		switch p.curTk.Kind {
		case token.MODULE:
			p.parsModule()
		}
	}
	return p.Ast
}

func (p *Parser) advance() *Parser {
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
		p.Tokens[p.pos] = tk
		p.pos++
		return p
	}
}

func (p *Parser) addString(s string) int {
	p.Ast.Strings[p.si] = s
	p.si++
	return p.si - 1
}

func (p *Parser) addNode(n ast.Node) int {
	p.Ast.Nodes[p.pn] = n
	p.pn++
	return p.pn - 1
}

func (p *Parser) match(kinds ...token.Kind) bool {
	for _, kind := range kinds {
		if p.curTk.Kind == kind {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) expect(kind token.Kind) bool {
	if p.curTk.Kind == kind {
		p.advance()
		return true
	} else {
		p.da.AddError(errors.Errors[5].IU("PARSER", []string{
			fmt.Sprintf("ожидалось: %s было получино: %s", kind.String(), p.curTk.Kind.String()),
		}))
		return false
	}
}

func (p *Parser) eat(kind token.Kind) bool {
	if p.curTk.Kind != kind {
		p.da.AddError(errors.Errors[5].IU("PARSER", []string{
			fmt.Sprintf("ожидалось: %s было получино: %s", kind.String(), p.curTk.Kind.String()),
		}))
		return false
	}
	p.advance()
	return true
}

func (p *Parser) VauleToken() string {
	return string(p.Tokens[p.pos].Literal(&p.input))
}
