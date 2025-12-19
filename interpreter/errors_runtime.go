package interpreter

import (
	"fmt"
)

type RuntimeError struct {
	*Position
	Msg string
}

func (e *RuntimeError) Error() string {
	return PrintError(e.Line, e.Column, "runtime", e.Msg)
}

// NewRuntimeError creates a generic runtime error
func NewRuntimeError(pos *Position, format string, args ...any) *RuntimeError {
	return &RuntimeError{
		Position: pos,
		Msg:      fmt.Sprintf(format, args...),
	}
}

func NewIdentifierNotFoundError(ident *Identifier) *RuntimeError {
	return &RuntimeError{
		Position: &Position{ident.Token.Line, ident.Token.Column},
		Msg:      fmt.Sprintf("identifier not found: %s", ident.Value),
	}
}

func NewVariableNotDefinedError(ident *Identifier) *RuntimeError {
	return &RuntimeError{
		Position: &Position{ident.Token.Line, ident.Token.Column},
		Msg:      fmt.Sprintf("variable not defined: %s", ident.Value),
	}
}

func NewDivisionByZeroError(pos *Position) *RuntimeError {
	return &RuntimeError{
		Position: pos,
		Msg:      "division by zero",
	}
}

func NewTypeMismatchError(pos *Position, expected, actual any) *RuntimeError {
	return &RuntimeError{
		Position: pos,
		Msg:      fmt.Sprintf("type mismatch: %T and %T", expected, actual),
	}
}

func NewUnknownOperatorError(pos *Position, op TokenType, left, right any) *RuntimeError {
	return &RuntimeError{
		Position: pos,
		Msg:      fmt.Sprintf("unknown operator or type: %v %s %v", left, op, right),
	}
}

func NewUnknownUnaryOperatorError(node *UnaryExpr, right any) *RuntimeError {
	return &RuntimeError{
		Position: &Position{node.Token.Line, node.Token.Column},
		Msg:      fmt.Sprintf("unknown unary operator: %s %v", node.Operator, right),
	}
}

func NewFunctionNotFoundError(node *CallExpr, name string) *RuntimeError {
	return &RuntimeError{
		Position: &Position{node.Token.Line, node.Token.Column},
		Msg:      fmt.Sprintf("function not found: %s", name),
	}
}

func NewInvalidFunctionCallError(node *CallExpr, reason string) *RuntimeError {
	return &RuntimeError{
		Position: &Position{node.Token.Line, node.Token.Column},
		Msg:      fmt.Sprintf("invalid function call: %s", reason),
	}
}

func NewInvalidRangeError(pos *Position, reason string) *RuntimeError {
	return &RuntimeError{
		Position: pos,
		Msg:      reason,
	}
}

func NewNegativeUintAssignmentError(pos *Position, value int64) *RuntimeError {
	return &RuntimeError{
		Position: pos,
		Msg:      fmt.Sprintf("cannot assign %d to uint: value must be non-negative", value),
	}
}

func NewInvalidUnofloatAssignmentError(pos *Position, value float64) *RuntimeError {
	return &RuntimeError{
		Position: pos,
		Msg:      fmt.Sprintf("cannot assign %f to unofloat: value out of range [0.0, 1.0]", value),
	}
}
