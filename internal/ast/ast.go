package ast

import "github.com/fugalang/fugu/internal/token"

type NodeType uint16
type Operator uint8

type AstArena struct {
	Nodes   []Node
	Strings []string
}

type Node struct {
	Type NodeType

	Data1 int
	Data2 int
	Data3 int
	Data4 float64
}

const (
	Invalid NodeType = iota
	Binary           // 1 + 2 типа
	Unary
	Literal

	Module
	Use

	If
)

const (
	OpInvalid Operator = iota
	OpPlus
	OpMinus
	OpMultiply
	OpDivide
	OpModulo
	OpPower
)

func Op(kind token.Kind) Operator {
	switch kind {
	case token.ADD:
		return OpPlus
	case token.SUB:
		return OpMinus
	case token.MUL:
		return OpMultiply
	case token.DIV:
		return OpDivide
	case token.MOD:
		return OpDivide
	case token.POW:
		return OpPower
	default:
		return OpInvalid
	}
}
