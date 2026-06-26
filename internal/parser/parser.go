package parser

import (
	"github.com/fugalang/fugu/internal/ast"
	"github.com/fugalang/fugu/internal/diagnostics"
	"github.com/fugalang/fugu/internal/lexer"
	"github.com/fugalang/fugu/internal/token"
)

type Parser struct {
	lex *lexer.Lexer

	Tokens []token.Token
	curTk  token.Token
	pastTk token.Token
	pos    int

	Ast *ast.AstArena
	da  *diagnostics.Arena
}

func New(input []byte, fileName string) *Parser {
	p := &Parser{
		Tokens: []token.Token{},

		Ast: &ast.AstArena{
			Nodes: []ast.Node{},
			Value: []ast.Value{},
		},

		da: &diagnostics.Arena{
			Source: string(input),
		},
	}
	p.lex = lexer.New(input, fileName, p.da)
	p.next()
	p.pos = 0
	return p
}

func (p *Parser) next() *Parser {
	tk := p.lex.NextToken()
	p.pastTk = p.curTk
	p.curTk = tk
	p.pos++

	if len(p.Tokens) > 0 && p.Tokens[len(p.Tokens)-1].Kind != token.EOF {
		p.Tokens = append(p.Tokens, tk)
	}

	return p
}

func (p *Parser) Parse() *ast.AstArena {
	switch p.curTk.Kind {
	case token.MODULE:
		p.module()
	}

	return p.Ast
}

func (p *Parser) addNode(n ast.Node) int {
	i := len(p.Ast.Nodes)
	p.Ast.Nodes = append(p.Ast.Nodes, n)
	return i
}

func (p *Parser) addValue(v ast.Value) int {
	i := len(p.Ast.Nodes)
	p.Ast.Value = append(p.Ast.Value, v)
	return i
}

func (p *Parser) VauleToken() string {
	return string(p.curTk.Literal(&p.lex.Input))
}

func (p *Parser) VaulePastToken() string {
	return string(p.pastTk.Literal(&p.lex.Input))
}

func (p *Parser) match(ks ...token.Kind) bool {
	for _, kind := range ks {
		if p.curTk.Kind == kind {
			p.next()
			return true
		}
	}
	return false
}
