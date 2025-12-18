package interpreter

import (
	"fmt"
)

type ParserError struct {
	*Position
	Msg string
}

func (e *ParserError) Error() string {
	return PrintError(e.Line, e.Column, "parser", e.Msg)
}

func NewParserError(pos *Position, format string, args ...any) *ParserError {
	return &ParserError{
		Position: pos,
		Msg:      fmt.Sprintf(format, args...),
	}
}

func NewIllegalTokenError(node *Token) *ParserError {
	return &ParserError{
		Position: &Position{
			Line:   node.Line,
			Column: node.Column,
		},
		Msg: fmt.Sprintf("Illegal token: %s", node.Literal),
	}
}

func NewExpectedTokenError(node *Token, expected TokenType) *ParserError {
	return &ParserError{
		Position: &Position{
			Line:   node.Line,
			Column: node.Column,
		},
		Msg: fmt.Sprintf("expected next token to be %s, got %s instead", expected, node.Type),
	}
}

func NewNoPrefixParseFnError(node *Token, t TokenType) *ParserError {
	return &ParserError{
		Position: &Position{
			Line:   node.Line,
			Column: node.Column,
		},
		Msg: fmt.Sprintf("no prefix parse function for %s found", t),
	}
}

func NewIntegerParseError(node *Token) *ParserError {
	return &ParserError{
		Position: &Position{
			Line:   node.Line,
			Column: node.Column,
		},
		Msg: fmt.Sprintf("could not parse %q as integer", node.Literal),
	}
}

func NewFloatParseError(node *Token) *ParserError {
	return &ParserError{
		Position: &Position{
			Line:   node.Line,
			Column: node.Column,
		},
		Msg: fmt.Sprintf("could not parse %q as float", node.Literal),
	}
}
