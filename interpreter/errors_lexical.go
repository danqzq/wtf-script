package interpreter

import "fmt"

type LexicalError struct {
	*Position
	Msg string
}

func (e *LexicalError) Error() string {
	return PrintError(e.Line, e.Column, "lexical", e.Msg)
}

func NewLexicalError(pos *Position, format string, args ...any) *LexicalError {
	return &LexicalError{
		Position: pos,
		Msg:      fmt.Sprintf(format, args...),
	}
}
