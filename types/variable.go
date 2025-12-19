package types

type Variable struct {
	Type  VarType
	Value any
}

type VarType int

// Distinct type for unofloat to differentiate from float64
type UnofloatType float64

const (
	Int VarType = iota
	Uint
	Float
	Unofloat // unofloat
	Bool
	String
	Unknown
)

func (t VarType) String() string {
	switch t {
	case Int:
		return "int"
	case Uint:
		return "uint"
	case Float:
		return "float"
	case Unofloat:
		return "unofloat"
	case Bool:
		return "bool"
	case String:
		return "string"
	default:
		return "unknown"
	}
}
