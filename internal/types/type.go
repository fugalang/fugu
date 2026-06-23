package types

type TypeKind uint8

const (
	Unknown TypeKind = iota

	Int // int64
	Int16
	Int32
	Int64

	Uint // Uint64
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

	Array
	Pointer
	Optional
	Channel
	Function
	Struct
	Interface
	Enum
	Generic
	None
)

type Type struct {
	Kind    TypeKind
	Dynamic bool
}
