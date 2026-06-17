package ast

type Arena struct {
	Nodes   []Node
	Strings []string
}

type StringID uint32
type NodeID uint32

const InvalidNode NodeID = ^NodeID(0)

type Node struct {
	Kind  NodeKind
	Flags uint16 // флаги: public, const, mut

	Data1 uint32
	Data2 uint32
	Data3 uint32
	Extra int64 // числовые значения
}

type NodeKind uint16

const (
	KindInvalid NodeKind = iota
	KindLiteral
	KindParenExpr
	KindAdditiveExpr
	KindMultiplicativeExpr
	KindPowerExpr
)

func (kind NodeKind) String() string {
	switch kind {
	case KindInvalid:
		return "KindInvalid"
	case KindLiteral:
		return "KindLiteral"
	case KindParenExpr:
		return "KindParenExpr"
	case KindAdditiveExpr:
		return "KindAdditiveExpr"
	case KindMultiplicativeExpr:
		return "KindMultiplicativeExpr"
	case KindPowerExpr:
		return "KindPowerExpr"
	default:
		return "KindInvalid"
	}
}
