package types

type TypeKind uint8

const (
	Unknown TypeKind = iota

	Int // int64
	Int16
	Int32
	int64

	Uint8
	Uint16
	Uint32
	Uint64

	Float32
	Float64

	Complex64
	Complex128

	Bool

	String
)

type Type struct {
	Kind TypeKind

	Array   bool
	Dynamic bool
}
