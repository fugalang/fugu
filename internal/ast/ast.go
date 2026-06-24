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

	Program
	Module
	Block

	Items
	Item
	Args
	ArgsTail

	Expr
	Primary
	AssignExpr
	BinaryExpr
	UnaryExpr
	CallExpr

	Literal
	Ident

	ExprStmt
	LetStmt
	UseStmt
	ReturnStmt
	IfStmt

	Assign
	Binary
	Unary
	Call
	Let
	Use
	Return
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

	case Program:
		return "Program"
	case Module:
		return "Module"
	case Block:
		return "Block"

	case Items:
		return "Items"
	case Item:
		return "Item"
	case Args:
		return "Args"
	case ArgsTail:
		return "ArgsTail"

	case Expr:
		return "Expr"
	case Primary:
		return "Primary"
	case AssignExpr:
		return "AssignExpr"
	case BinaryExpr:
		return "BinaryExpr"
	case UnaryExpr:
		return "UnaryExpr"
	case CallExpr:
		return "CallExpr"

	case Literal:
		return "Literal"
	case Ident:
		return "Ident"

	case ExprStmt:
		return "ExprStmt"
	case LetStmt:
		return "LetStmt"
	case UseStmt:
		return "UseStmt"
	case ReturnStmt:
		return "ReturnStmt"
	case IfStmt:
		return "IfStmt"

	case Assign:
		return "Assign"
	case Binary:
		return "Binary"
	case Unary:
		return "Unary"
	case Call:
		return "Call"
	case Let:
		return "Let"
	case Use:
		return "Use"
	case Return:
		return "Return"
	case If:
		return "If"

	default:
		return "Unknown"
	}
}
