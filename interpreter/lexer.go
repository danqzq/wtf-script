package interpreter

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Lexer implements a concurrent lexical shi as a scanner using state functions and  channels
// Inspired by Rob Pike's "Lexical Scanning in Go"

type Lexer struct {
	name        string
	input       string
	start       int
	pos         int
	width       int
	line        int // 1 - indexed
	column      int // 1 - indexed
	startLine   int // start line of current token
	startColumn int // start column of current token
	tokens      chan Token
	state       stateFn
}

type stateFn func(*Lexer) stateFn

const EOS = -1 // end of shi

func NewLexer(name, input string) *Lexer {
	l := &Lexer{
		name:        name,
		input:       input,
		tokens:      make(chan Token),
		line:        1,
		column:      1,
		startLine:   1,
		startColumn: 1,
	}
	go l.run() // run the state machine concurrently
	return l
}

// run executes the state machine
func (l *Lexer) run() {
	for l.state = lexStart; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.tokens)
}

func (l *Lexer) NextToken() Token {
	return <-l.tokens
}

// sends the token back to the client
func (l *Lexer) emit(t TokenType) {
	l.tokens <- Token{
		Type:    t,
		Literal: l.input[l.start:l.pos],
		Line:    l.startLine,
		Column:  l.startColumn,
	}
	l.start = l.pos
	l.startLine = l.line
	l.startColumn = l.column
}

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return EOS
	}
	inputCurrPos := l.input[l.pos:]
	ch, width := utf8.DecodeRuneInString(inputCurrPos)
	l.width = width
	l.pos += l.width

	if ch == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}

	return ch
}

func (l *Lexer) peek() rune {
	ch := l.next()
	l.backup()
	return ch
}

func (l *Lexer) backup() {
	l.pos -= l.width
	if l.pos >= 0 && l.pos < len(l.input) && l.input[l.pos] == '\n' {
		l.line--
		col := 1
		for i := l.pos - 1; i >= 0; i-- {
			if l.input[i] == '\n' {
				break
			}
			col++
		}
		l.column = col
	} else {
		l.column--
	}
}

func (l *Lexer) ignore() {
	l.start = l.pos
	l.startLine = l.line
	l.startColumn = l.column
}

func (l *Lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set
func (l *Lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func (l *Lexer) errorf(formattedMsg string, args ...any) stateFn {
	l.tokens <- Token{
		Type:    ILLEGAL,
		Literal: fmt.Sprintf(formattedMsg, args...),
		Line:    l.line,
		Column:  l.column,
	}
	return lexStart
}

func lexStart(l *Lexer) stateFn {
	for {
		ch := l.next()

		switch {
		case ch == EOS:
			l.emit(EOF)
			return nil
		case isSpace(ch):
			l.ignore()
		case ch == '/':
			if l.peek() == '/' {
				return lexComment
			}
			l.emit(SLASH)
		case ch == '+':
			l.emit(PLUS)
		case ch == '-':
			// Could be minus or start of negative number
			if isDigit(l.peek()) {
				l.backup()
				return lexNumber
			}
			l.emit(MINUS)
		case ch == '*':
			l.emit(ASTERISK)
		case ch == '=':
			l.emit(ASSIGN)
		case ch == ';':
			l.emit(SEMICOLON)
		case ch == ',':
			l.emit(COMMA)
		case ch == '(':
			l.emit(LPAREN)
		case ch == ')':
			l.emit(RPAREN)
		case ch == '"':
			return lexString
		case isDigit(ch):
			l.backup()
			return lexNumber
		case isAlpha(ch):
			l.backup()
			return lexIdentifier
		default:
			return l.errorf("unexpected character: %c", ch)
		}
	}
}

func lexComment(l *Lexer) stateFn {
	l.next() // consume second '/'
	for {
		ch := l.next()
		if ch == '\n' || ch == EOS {
			l.backup()
			l.ignore() // Don't emit comments as tokens
			return lexStart
		}
	}
}

func lexString(l *Lexer) stateFn {
	for {
		ch := l.next()
		switch ch {
		case EOS, '\n':
			l.emit(STRING)
			return l.errorf("unterminated string")
		case '"':
			l.emit(STRING)
			return lexStart
		case '\\':
			// Handle escape sequences if needed
			if l.next() == EOS {
				return l.errorf("unterminated string")
			}
		}
	}
}

func lexNumber(l *Lexer) stateFn {
	l.accept("+-")
	digits := "0123456789"
	l.acceptRun(digits)

	// Check for decimal point
	if l.accept(".") {
		l.acceptRun(digits)
		l.emit(FLOAT)
	} else {
		l.emit(INT)
	}
	return lexStart
}

func lexIdentifier(l *Lexer) stateFn {
	for {
		ch := l.next()
		if !isAlphaNumeric(ch) {
			l.backup()
			break
		}
	}

	word := l.input[l.start:l.pos]
	tokType := LookupIdent(word)
	l.emit(tokType)

	return lexStart
}

func isSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

func isAlpha(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isAlphaNumeric(ch rune) bool {
	return isAlpha(ch) || isDigit(ch)
}
