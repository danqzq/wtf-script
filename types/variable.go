package types

import "fmt"

type Variable struct {
	Type  VarType
	Value any
}

type VarType int

const (
	Int VarType = iota
	Uint
	Float
	UnoFloat
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
	case UnoFloat:
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
	case Float, UnoFloat:
		return fmt.Sprintf("%f", v.Value.(float64))
	case Bool:
		return fmt.Sprintf("%t", v.Value.(bool))
	case String:
		return v.Value.(string)
	default:
		return fmt.Sprintf("%v", v.Value)
	}
}
