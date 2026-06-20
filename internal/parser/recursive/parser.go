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
	peekTk token.Token
	pos    int

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
	p.advance().advance()
	return p
}

func (p *Parser) Parse() *ast.AstArena {
	for p.curTk.Kind != token.EOF {
		fmt.Println(p.curTk.Kind.String())
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

		if tk.Kind == token.EOF && p.peekTk.Kind == token.EOF {
			p.curTk = p.peekTk
			return p
		}

		p.curTk = p.peekTk
		p.peekTk = tk
		p.pos++
		p.Tokens = append(p.Tokens, tk)
		return p
	}
}

func (p *Parser) match(kind token.Kind) bool {
	if p.curTk.Kind == kind {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) matchPeek(kind token.Kind) bool {
	return p.peekTk.Kind == kind
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
