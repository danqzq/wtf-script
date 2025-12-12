package interpreter

// token.go  - Token definitions
// we only handle the tokens we need right now

type TokenType string

type Token struct {
	Type    TokenType
	Literal string // the actual text
	Line    int    // line number for error messages
	Column  int    // Column number for error messages
}

const (
	// Special tokens
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers and literals
	IDENT TokenType = "IDENT"
	INT   TokenType = "INT"

	// Operators
	ASSIGN TokenType = "="
	PLUS   TokenType = "+"
	MINUS  TokenType = "-"

	// Delimiters
	SEMICOLON TokenType = ";"
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"

	// Keywords  with one type now
	TYPE_INT TokenType = "INT_TYPE"
)

// keywords maps keyword strings to their TokenType
var keywords = map[string]TokenType{
	"int": TYPE_INT,
}

func LookupIdent(ident string) TokenType {
	if token, ok := keywords[ident]; ok {
		return token
	}
	return IDENT
}
