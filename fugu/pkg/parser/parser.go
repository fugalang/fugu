package parser

import (
	"fugu/pkg/diagnostics"
	lexer "fugu/pkg/lexer"
	"fugu/pkg/parser/action"
	"fugu/pkg/parser/action/ast"
	"fugu/pkg/token"
)

type Parser struct {
	Tokens []token.Token
	Roots  []ast.NodeID

	lex   *lexer.Lexer
	diagn *diagnostics.Diagnostics

	curToken token.Token
	ast      *ast.Arena

	pos int
}

func New(input []byte, fileName string) *Parser {
	lex := lexer.New(input, fileName)

	p := &Parser{
		Tokens: make([]token.Token, 0),
		Roots:  make([]ast.NodeID, 0),

		lex: lex,
		ast: &ast.Arena{
			Nodes:   make([]ast.Node, 0, 512),
			Strings: make([]string, 0, 256),
		},
	}

	p.diagn = p.lex.Report()

	if p.diagn.IsUse {
		p.diagn.SendTk(diagnostics.ParserCantStartWork, p.curToken)
		return p
	}

	p.advance()

	if p.curToken.Kind == token.EOF {
		p.diagn.SendTk(diagnostics.ParserCantStartWork, p.curToken)
		return p
	}

	return p
}

func (p *Parser) advance() {
	p.curToken = p.lex.NextToken()

	if len(p.Tokens) == 0 || p.Tokens[len(p.Tokens)-1].Kind != token.EOF {
		p.Tokens = append(p.Tokens, p.curToken)
	}

	p.pos++
}

type stackSymbol struct {
	isNode bool
	node   ast.NodeID
	tok    token.Token
}

func (p *Parser) Run() {
	stateStack := []int{0}
	symStack := []stackSymbol{}

	for {

		if p.skip() {
			continue
		}

		state := stateStack[len(stateStack)-1]

		act := action.Action(state, p.curToken.Kind)

		switch act.Typ {
		case action.Accept:
			if len(symStack) > 0 {
				last := symStack[len(symStack)-1]
				if last.isNode {
					p.Roots = append(p.Roots, last.node)
				}
			}
			return

		case action.Error:
			p.diagn.SendTk(act.ErrCode, p.curToken)
			return

		case action.Shift:
			symStack = append(symStack, stackSymbol{
				isNode: false,
				tok:    p.curToken,
			})

			stateStack = append(stateStack, act.State)
			p.advance()

		case action.Reduce:
			count := act.State

			if count > len(symStack) {
				p.diagn.SendTk(diagnostics.ParserCantStartWork, p.curToken)
				return
			}

			start := len(symStack) - count
			items := symStack[start:]
			symStack = symStack[:start]

			// build AST node
			newNode := p.buildNode(act.NodeKind, items)
			if newNode == ast.InvalidNode {
				p.diagn.SendTk(diagnostics.ParserCantStartWork, p.curToken)
				return
			}

			stateStack = stateStack[:len(stateStack)-count]
			top := stateStack[len(stateStack)-1]
			next := action.Goto(top, act.NodeKind)

			if next < 0 {
				p.diagn.SendTk(diagnostics.ParserCantStartWork, p.curToken)
				return
			}

			stateStack = append(stateStack, next)
			symStack = append(symStack, stackSymbol{
				isNode: true,
				node:   newNode,
			})
		}
	}
}

type nodeBuilder func(*Parser, ast.NodeKind, []stackSymbol) ast.NodeID

var nodeBuilders = map[ast.NodeKind]nodeBuilder{
	ast.KindLiteral:            buildLiteral,
	ast.KindAdditiveExpr:       buildBinary,
	ast.KindMultiplicativeExpr: buildBinary,
	ast.KindDegreeExpr:         buildBinary,
}

func (p *Parser) buildNode(kind ast.NodeKind, items []stackSymbol) ast.NodeID {
	if builder, ok := nodeBuilders[kind]; ok {
		return builder(p, kind, items)
	}

	p.diagn.SendTk(diagnostics.StateDoesNotToken, p.curToken)
	return ast.InvalidNode
}

func buildLiteral(p *Parser, kind ast.NodeKind, items []stackSymbol) ast.NodeID {
	if len(items) == 3 {
		return p.ensureNode(items[1])
	}
	if len(items) != 1 {
		return ast.InvalidNode
	}

	if items[0].isNode {
		return items[0].node
	}

	return p.makeLiteral(items[0].tok, kind)
}

func buildBinary(p *Parser, kind ast.NodeKind, items []stackSymbol) ast.NodeID {
	if len(items) == 1 {
		return p.ensureNode(items[0])
	}
	if len(items) != 3 {
		return ast.InvalidNode
	}

	left := p.ensureNode(items[0])
	right := p.ensureNode(items[2])

	id := ast.NodeID(len(p.ast.Nodes))

	p.ast.Nodes = append(p.ast.Nodes, ast.Node{
		Kind:  kind,
		Data1: uint32(left),
		Data2: uint32(right),
		Data3: uint32(items[1].tok.Kind),
	})

	return id
}

func (p *Parser) ensureNode(s stackSymbol) ast.NodeID {
	if s.isNode {
		return s.node
	}

	return p.makeLiteral(s.tok, ast.KindLiteral)
}

func (p *Parser) makeLiteral(tk token.Token, kind ast.NodeKind) ast.NodeID {
	id := ast.NodeID(len(p.ast.Nodes))

	p.ast.Nodes = append(p.ast.Nodes, ast.Node{
		Kind:  kind,
		Data1: uint32(tk.Kind),
		Data2: uint32(tk.Start),
		Data3: uint32(tk.End),
	})

	return id
}

func (p *Parser) skip() bool {
	for {
		switch p.curToken.Kind {
		case token.SPACING, token.COMMENT, token.M_COMMENT:
			p.advance()
			return true
		default:
			return false
		}
	}
}
