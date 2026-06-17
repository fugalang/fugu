package parser

import (
	lexer "fugu/pkg/lexer"
	"fugu/pkg/parser/action/ast"
	"fugu/pkg/reporter"
	"fugu/pkg/token"
)

type Parser struct {
	Tokens []token.Token
	Roots  []ast.NodeID

	lex    *lexer.Lexer
	report *reporter.Reporter

	curToken token.Token
	ast      *ast.Arena

	pos int
}

func New(input []byte, fileName string) *Parser {
	lex := lexer.New(input, fileName)

	pars := &Parser{
		Tokens: make([]token.Token, 0),
		Roots:  make([]ast.NodeID, 0),

		lex: lex,

		ast: &ast.Arena{
			Nodes:   make([]ast.Node, 0, 512),
			Strings: make([]string, 0, 256),
		},
	}
	pars.report = pars.lex.Report()
	if pars.report.IsUse {
		pars.report.SendTk(reporter.ParserCantStartWork, pars.curToken)
		return pars
	}

	pars.advance()
	if pars.curToken.Kind == token.EOF {
		pars.report.SendTk(reporter.ParserCantStartWork, pars.curToken)
		return pars
	}

	return pars
}

func (ps *Parser) advance() *Parser {
	ps.curToken = ps.lex.NextToken()

	if len(ps.Tokens) == 0 || ps.Tokens[len(ps.Tokens)-1].Kind != token.EOF {
		ps.Tokens = append(ps.Tokens, ps.curToken)
	}
	ps.pos++

	return ps
}

func (ps *Parser) Run() {
	for ps.curToken.Kind != token.EOF {
		ps.skipSlag()
		if ps.curToken.Kind == token.EOF {
			return
		}

		pos := ps.pos
		node := ps.top()

		if node != ast.InvalidNode {
			ps.Roots = append(ps.Roots, node)
		}

		ps.endTop()

		if ps.pos == pos {
			ps.advance()
		}
	}
}

func (ps *Parser) top() ast.NodeID {
	switch ps.curToken.Kind {
	case token.MODULE, token.USE, token.FN, token.LET, token.CONST, token.TYPE, token.ENUM, token.STRUCT, token.INTERFACE:
		return ps.decl()
	case token.RETURN, token.IF, token.MATCH, token.FOR, token.CONTINUE, token.BREAK:
		return ps.stmt()
	default:
		return ps.expr(1)
	}
}

func (ps *Parser) decl() ast.NodeID {
	switch ps.curToken.Kind {
	case token.MODULE:
		return ps.moduleDecl()
	case token.USE:
		return ps.useDecl()
	case token.FN:
		return ps.fnDecl()
	case token.LET, token.CONST:
		return ps.varDecl()
	case token.TYPE, token.ENUM, token.STRUCT, token.INTERFACE:
		return ps.typeDecl()
	default:
		ps.report.SendTk(reporter.StateDoesNotToken, ps.curToken)
		return ast.InvalidNode
	}
}

func (ps *Parser) stmt() ast.NodeID {
	switch ps.curToken.Kind {
	default:
		ps.report.SendTk(reporter.StateDoesNotToken, ps.curToken)
		return ast.InvalidNode
	}
}

func (ps *Parser) moduleDecl() ast.NodeID {
	ps.report.SendTk(reporter.StateDoesNotToken, ps.curToken)
	return ast.InvalidNode
}

func (ps *Parser) useDecl() ast.NodeID {
	ps.report.SendTk(reporter.StateDoesNotToken, ps.curToken)
	return ast.InvalidNode
}

func (ps *Parser) fnDecl() ast.NodeID {
	ps.report.SendTk(reporter.StateDoesNotToken, ps.curToken)
	return ast.InvalidNode
}

func (ps *Parser) varDecl() ast.NodeID {
	ps.report.SendTk(reporter.StateDoesNotToken, ps.curToken)
	return ast.InvalidNode
}

func (ps *Parser) typeDecl() ast.NodeID {
	ps.report.SendTk(reporter.StateDoesNotToken, ps.curToken)
	return ast.InvalidNode
}

func (ps *Parser) expr(minPrior int) ast.NodeID {
	left := ps.term()

	for {
		ps.skipSlag()

		prior, right, ok := opPrior(ps.curToken.Kind)
		if !ok || prior < minPrior {
			return left
		}

		op := ps.curToken
		ps.advance()

		next := prior + 1
		if right {
			next = prior
		}

		rightNode := ps.expr(next)
		left = ps.node(ast.Node{
			Kind:  opKind(op.Kind),
			Data1: uint32(left),
			Data2: uint32(rightNode),
			Data3: uint32(op.Kind),
		})
	}
}

func (ps *Parser) term() ast.NodeID {
	ps.skipSlag()

	switch ps.curToken.Kind.Group() {
	case token.GNUMBER, token.GSTRING, token.GLITERAL:
		tk := ps.curToken
		ps.advance()
		return ps.node(ast.Node{
			Kind:  ast.KindLiteral,
			Data1: uint32(tk.Kind),
			Data2: uint32(tk.Start),
			Data3: uint32(tk.End),
		})
	}

	if ps.curToken.Kind == token.L_PAREN {
		ps.advance()
		node := ps.expr(1)
		ps.skipSlag()

		if ps.curToken.Kind != token.R_PAREN {
			ps.report.SendTk(reporter.StateDoesNotToken, ps.curToken)
			return node
		}

		ps.advance()
		return ps.node(ast.Node{
			Kind:  ast.KindParenExpr,
			Data1: uint32(node),
		})
	}

	ps.report.SendTk(reporter.StateDoesNotToken, ps.curToken)
	if ps.curToken.Kind != token.EOF {
		ps.advance()
	}
	return ps.node(ast.Node{Kind: ast.KindInvalid})
}

func (ps *Parser) skipSlag() {
	for ps.curToken.Kind == token.SPACING || ps.curToken.Kind == token.COMMENT || ps.curToken.Kind == token.M_COMMENT {
		ps.advance()
	}
}

func (ps *Parser) node(node ast.Node) ast.NodeID {
	ps.ast.Nodes = append(ps.ast.Nodes, node)
	return ast.NodeID(len(ps.ast.Nodes) - 1)
}

func (ps *Parser) endTop() {
	ps.skipSlag()
	if ps.curToken.Kind == token.END {
		ps.advance()
	}
}

func opPrior(kind token.TokenKind) (prior int, right bool, ok bool) {
	switch kind {
	case token.INCREASE, token.DECREASE:
		return 1, false, true
	case token.MULTIPLY, token.DIVIDE, token.REMAINDER:
		return 2, false, true
	case token.DEGREE:
		return 3, true, true
	default:
		return 0, false, false
	}
}

func opKind(kind token.TokenKind) ast.NodeKind {
	switch kind {
	case token.INCREASE, token.DECREASE:
		return ast.KindAdditiveExpr
	case token.MULTIPLY, token.DIVIDE, token.REMAINDER:
		return ast.KindMultiplicativeExpr
	case token.DEGREE:
		return ast.KindPowerExpr
	default:
		return ast.KindInvalid
	}
}
