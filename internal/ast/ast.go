package ast

type NodeType uint16

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
