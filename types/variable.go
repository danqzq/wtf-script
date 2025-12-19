package types

type Variable struct {
	Type  VarType
	Value any
}

type VarType int

// Distinct type for unofloat to differentiate from float64
type Unofloat float64

const (
	Int VarType = iota
	Uint
	Float
	UnitFloat // unofloat
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
	case UnitFloat:
		return "unofloat"
	case Bool:
		return "bool"
	case String:
		return "string"
	default:
		return "unknown"
	}
}
