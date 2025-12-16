package interpreter

import "fmt"

type LexicalError struct {
	Line   int
	Column int
	Msg    string
}

func (e *LexicalError) Error() string {
	return PrintError(e.Line, e.Column, "lexical", e.Msg)
}

func NewLexicalError(line, col int, format string, args ...any) *LexicalError {
	return &LexicalError{
		Line:   line,
		Column: col,
		Msg:    fmt.Sprintf(format, args...),
	}
}
