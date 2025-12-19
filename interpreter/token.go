package interpreter

// Token definitions for WTFScript
type TokenType string

type Token struct {
	Type    TokenType
	Literal string // the actual text
	Line    int    // line number for error messages
	Column  int    // column number for error messages
}

const (
	// Special tokens
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"
	COMMENT TokenType = "COMMENT"

	// Identifiers and literals
	IDENT  TokenType = "IDENT"
	INT    TokenType = "INT"
	FLOAT  TokenType = "FLOAT"
	STRING TokenType = "STRING"
	TRUE   TokenType = "TRUE"
	FALSE  TokenType = "FALSE"

	// Operators
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"

	// Comparison operators
	EQ   TokenType = "=="
	NEQ  TokenType = "!="
	LT   TokenType = "<"
	LTE  TokenType = "<="
	GT   TokenType = ">"
	GTE  TokenType = ">="
	BANG TokenType = "!"

	// Logical operators
	AND TokenType = "&&"
	OR  TokenType = "||"

	// Delimiters
	SEMICOLON TokenType = ";"
	COMMA     TokenType = ","
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"
	LBRACE    TokenType = "{"
	RBRACE    TokenType = "}"

	// Type keywords
	TYPE_INT      TokenType = "INT_TYPE"
	TYPE_UINT     TokenType = "UINT_TYPE"
	TYPE_FLOAT    TokenType = "FLOAT_TYPE"
	TYPE_UNOFLOAT TokenType = "UNOFLOAT_TYPE"
	TYPE_BOOL     TokenType = "BOOL_TYPE"
	TYPE_STRING   TokenType = "STRING_TYPE"

	// Control flow keywords
	IF     TokenType = "IF"
	ELSE   TokenType = "ELSE"
	IFRAND TokenType = "IFRAND"
)

// keywords maps keyword strings to their TokenType
var keywords = map[string]TokenType{
	"int":      TYPE_INT,
	"uint":     TYPE_UINT,
	"float":    TYPE_FLOAT,
	"unofloat": TYPE_UNOFLOAT,
	"bool":     TYPE_BOOL,
	"string":   TYPE_STRING,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"ifrand":   IFRAND,
}

// LookupIdent checks if an identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// String returns a human-readable representation of the token
func (t Token) String() string {
	return string(t.Type) + ":" + t.Literal
}
