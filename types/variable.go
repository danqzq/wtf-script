package types

import "fmt"

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

func (v Variable) ToString() string {
	switch v.Type {
	case Int:
		return fmt.Sprintf("%d", v.Value.(int))
	case Uint:
		return fmt.Sprintf("%d", v.Value.(uint64))
	case Float:
		return fmt.Sprintf("%f", v.Value.(float64))
	case UnitFloat:
		return fmt.Sprintf("%f", v.Value.(Unofloat))
	case Bool:
		return fmt.Sprintf("%t", v.Value.(bool))
	case String:
		return v.Value.(string)
	default:
		return fmt.Sprintf("%v", v.Value)
	}
}
