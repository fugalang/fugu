package ast

import (
	"github.com/fugalang/fugu/internal/token"
)

type NodeType uint16
type Operator uint8
type ValueKind uint8

type AstArena struct {
	Nodes []Node
	Value []Value
}

type Node struct {
	Type NodeType

	Data1 int
	Data2 int
	Data3 int
	Data4 float64
}

type Value struct {
	Type ValueKind

	STR  string
	I64  int64
	F64  float64
	C128 complex128
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
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMod
	OpPow
)

const (
	_ ValueKind = iota
	String
	Int
	Float
	Complex
)

func Op(kind token.Kind) Operator {
	switch kind {
	case token.ADD:
		return OpAdd
	case token.SUB:
		return OpSub
	case token.MUL:
		return OpMul
	case token.DIV:
		return OpDiv
	case token.MOD:
		return OpMod
	case token.POW:
		return OpPow
	default:
		return OpInvalid
	}
}
