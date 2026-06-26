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

type TemplatePart struct {
	IsExpr bool
	Str    string
	Expr   int
}

type Value struct {
	Type ValueKind

	S8   string
	I64  int64
	F64  float64
	C128 complex128
}

const (
	_ ValueKind = iota
	String
	Int
	Float
	Complex
	Char
)

const (
	Invalid NodeType = iota

	BinaryExpr
	UnaryExpr

	Literal
	Ident

	Template

	ModuleDecl
)

const (
	OpInvalid Operator = iota

	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMod
	OpPow

	OpEq
	OpNeq
	OpLt
	OpGt
	OpLe
	OpGe

	OpAnd
	OpOr

	OpNeg
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

	case token.EQ:
		return OpEq
	case token.NEQ:
		return OpNeq
	case token.LT:
		return OpLt
	case token.GT:
		return OpGt
	case token.LE:
		return OpLe
	case token.GE:
		return OpGe

	case token.AND:
		return OpAnd
	case token.OR:
		return OpOr

	default:
		return OpInvalid
	}
}

func (n NodeType) String() string {
	switch n {
	case Invalid:
		return "Invalid"

	default:
		return "Unknown"
	}
}
