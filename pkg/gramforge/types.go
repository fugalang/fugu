package gram

import "fmt"

const EOF = "EOF"

type Assoc int

const (
	AssocNone Assoc = iota
	AssocLeft
	AssocRight
	AssocNonassoc
)

type Grammar struct {
	Start       string
	Starts      []string
	Terminals   map[string]bool
	terminalSeq []string
	precedence  map[string]Prec
	productions []Production
	imports     map[string]bool
}

type Prec struct {
	Level int
	Assoc Assoc
}

type Production struct {
	ID     int
	LHS    string
	RHS    []string
	Node   string
	PrecBy string
}

type Table struct {
	Grammar *Grammar
	States  []State
	Actions map[int]map[string]Action
	Gotos   map[int]map[string]int
}

type State struct {
	Items []Item
}

type Item struct {
	Prod int
	Dot  int
}

type ActionKind int

const (
	ActionShift ActionKind = iota + 1
	ActionReduce
	ActionAccept
)

type Action struct {
	Kind  ActionKind
	State int
	Prod  Production
}

func New() *Grammar {
	return &Grammar{
		Terminals:  make(map[string]bool),
		precedence: make(map[string]Prec),
		Starts:     make([]string, 0),
		imports:    make(map[string]bool),
	}
}

func (g *Grammar) AddToken(name string) {
	if name == "" || g.Terminals[name] {
		return
	}
	g.Terminals[name] = true
	g.terminalSeq = append(g.terminalSeq, name)
}

func (g *Grammar) AddPrecedence(assoc Assoc, tokens ...string) {
	level := 1
	for _, p := range g.precedence {
		if p.Level >= level {
			level = p.Level + 1
		}
	}

	for _, tk := range tokens {
		g.AddToken(tk)
		g.precedence[tk] = Prec{Level: level, Assoc: assoc}
	}
}

func (g *Grammar) AddProduction(lhs string, rhs []string, node string, precBy string) error {
	if lhs == "" {
		return fmt.Errorf("грамматика: пустая левая часть правила")
	}
	if node == "" {
		node = lhs
	}

	p := Production{
		ID:     len(g.productions),
		LHS:    lhs,
		RHS:    append([]string(nil), rhs...),
		Node:   node,
		PrecBy: precBy,
	}
	g.productions = append(g.productions, p)

	if g.Start == "" && len(g.Starts) == 0 {
		g.Start = lhs
	}
	return nil
}

func (g *Grammar) completeTerminals() {
	nonterms := g.nonTerminals()
	for _, p := range g.productions {
		for _, sym := range p.RHS {
			if !nonterms[sym] {
				g.AddToken(sym)
			}
		}
		if p.PrecBy != "" {
			g.AddToken(p.PrecBy)
		}
	}
	g.AddToken(EOF)
}

func (g *Grammar) nonTerminals() map[string]bool {
	out := make(map[string]bool)
	for _, p := range g.productions {
		out[p.LHS] = true
	}
	return out
}

func (g *Grammar) isNonTerminal(sym string) bool {
	for _, p := range g.productions {
		if p.LHS == sym {
			return true
		}
	}
	return false
}

func (g *Grammar) isTerminal(sym string) bool {
	return g.Terminals[sym] || !g.isNonTerminal(sym)
}

func (g *Grammar) GetPrimaryStart() string {
	if len(g.Starts) > 0 {
		return g.Starts[0]
	}
	return g.Start
}

func (g *Grammar) AddStart(name string) error {
	if name == "" {
		return fmt.Errorf("грамматика: пустой стартовый символ")
	}
	for _, s := range g.Starts {
		if s == name {
			return nil
		}
	}
	g.Starts = append(g.Starts, name)
	if g.Start == "" {
		g.Start = name
	}
	return nil
}
