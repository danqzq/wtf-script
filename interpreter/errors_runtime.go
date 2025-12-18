package interpreter

import (
	"fmt"
)

type RuntimeError struct {
	Line   int
	Column int
	Msg    string
}

func (e *RuntimeError) Error() string {
	return PrintError(e.Line, e.Column, "runtime", e.Msg)
}

// NewRuntimeError creates a generic runtime error
func NewRuntimeError(line, col int, format string, args ...any) *RuntimeError {
	return &RuntimeError{
		Line:   line,
		Column: col,
		Msg:    fmt.Sprintf(format, args...),
	}
}

// Specific Error Constructors

func NewIdentifierNotFoundError(ident *Identifier) *RuntimeError {
	return &RuntimeError{
		Line:   ident.Token.Line,
		Column: ident.Token.Column,
		Msg:    fmt.Sprintf("identifier not found: %s", ident.Value),
	}
}

func NewVariableNotDefinedError(ident *Identifier) *RuntimeError {
	return &RuntimeError{
		Line:   ident.Token.Line,
		Column: ident.Token.Column,
		Msg:    fmt.Sprintf("variable not defined: %s", ident.Value),
	}
}

func NewDivisionByZeroError(line, col int) *RuntimeError {
	return &RuntimeError{
		Line:   line,
		Column: col,
		Msg:    "division by zero",
	}
}

func NewTypeMismatchError(line, col int, expected, actual any) *RuntimeError {
	return &RuntimeError{
		Line:   line,
		Column: col,
		Msg:    fmt.Sprintf("type mismatch: %T and %T", expected, actual),
	}
}

func NewUnknownOperatorError(line, col int, op string, left, right any) *RuntimeError {
	return &RuntimeError{
		Line:   line,
		Column: col,
		Msg:    fmt.Sprintf("unknown operator or type: %v %s %v", left, op, right),
	}
}

func NewUnknownUnaryOperatorError(node *UnaryExpr, right any) *RuntimeError {
	return &RuntimeError{
		Line:   node.Token.Line,
		Column: node.Token.Column,
		Msg:    fmt.Sprintf("unknown unary operator: %s %v", node.Operator, right),
	}
}

func NewFunctionNotFoundError(node *CallExpr, name string) *RuntimeError {
	return &RuntimeError{
		Line:   node.Token.Line,
		Column: node.Token.Column,
		Msg:    fmt.Sprintf("function not found: %s", name),
	}
}

func NewInvalidFunctionCallError(node *CallExpr, reason string) *RuntimeError {
	return &RuntimeError{
		Line:   node.Token.Line,
		Column: node.Token.Column,
		Msg:    fmt.Sprintf("invalid function call: %s", reason),
	}
}

func NewInvalidRangeError(line, col int, reason string) *RuntimeError {
	return &RuntimeError{
		Line:   line,
		Column: col,
		Msg:    reason,
	}
}

func NewNegativeUintAssignmentError(line, col int, value int64) *RuntimeError {
	return &RuntimeError{
		Line:   line,
		Column: col,
		Msg:    fmt.Sprintf("cannot assign %d to uint: value must be non-negative", value),
	}
}

func NewInvalidUnofloatAssignmentError(line, col int, value float64) *RuntimeError {
	return &RuntimeError{
		Line:   line,
		Column: col,
		Msg:    fmt.Sprintf("cannot assign %f to unofloat: value out of range [0.0, 1.0]", value),
	}
}
