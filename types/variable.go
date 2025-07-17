package types

import "fmt"

type Variable struct {
	Type  VarType
	Value interface{}
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
