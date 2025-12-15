package interpreter

import (
	"testing"
)

func TestLexer_BasicTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []TokenType
	}{
		{
			name:     "simple variable declaration",
			input:    "int x;",
			expected: []TokenType{TYPE_INT, IDENT, SEMICOLON, EOF},
		},
		{
			name:     "variable assignment",
			input:    "int x = 5;",
			expected: []TokenType{TYPE_INT, IDENT, ASSIGN, INT, SEMICOLON, EOF},
		},
		{
			name:     "arithmetic expression",
			input:    "int y = 10 + 20 * 3;",
			expected: []TokenType{TYPE_INT, IDENT, ASSIGN, INT, PLUS, INT, ASTERISK, INT, SEMICOLON, EOF},
		},
		{
			name:     "function call",
			input:    "print(x);",
			expected: []TokenType{IDENT, LPAREN, IDENT, RPAREN, SEMICOLON, EOF},
		},
		{
			name:     "string literal",
			input:    `string s = "Hello World";`,
			expected: []TokenType{TYPE_STRING, IDENT, ASSIGN, STRING, SEMICOLON, EOF},
		},
		{
			name:     "boolean values",
			input:    "bool flag = true;",
			expected: []TokenType{TYPE_BOOL, IDENT, ASSIGN, TRUE, SEMICOLON, EOF},
		},
		{
			name:     "float number",
			input:    "float pi = 3.14;",
			expected: []TokenType{TYPE_FLOAT, IDENT, ASSIGN, FLOAT, SEMICOLON, EOF},
		},
		{
			name:     "range declaration",
			input:    "int(0, 100) x;",
			expected: []TokenType{TYPE_INT, LPAREN, INT, COMMA, INT, RPAREN, IDENT, SEMICOLON, EOF},
		},
		{
			name:     "comment handling",
			input:    "// This is a comment\nint x;",
			expected: []TokenType{TYPE_INT, IDENT, SEMICOLON, EOF},
		},
		{
			name:     "all type keywords",
			input:    "int uint float unofloat bool string",
			expected: []TokenType{TYPE_INT, TYPE_UINT, TYPE_FLOAT, TYPE_UNOFLOAT, TYPE_BOOL, TYPE_STRING, EOF},
		},
		{
			name:     "negative number",
			input:    "int x = -5;",
			expected: []TokenType{TYPE_INT, IDENT, ASSIGN, INT, SEMICOLON, EOF},
		},
		{
			name:     "parenthesized expression",
			input:    "int result = (10 + 5) * 2;",
			expected: []TokenType{TYPE_INT, IDENT, ASSIGN, LPAREN, INT, PLUS, INT, RPAREN, ASTERISK, INT, SEMICOLON, EOF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.name, tt.input)

			for i, expectedType := range tt.expected {
				tok := lexer.NextToken()
				if tok.Type != expectedType {
					t.Errorf("test[%d] - token type wrong. expected=%q, got=%q (literal: %q)",
						i, expectedType, tok.Type, tok.Literal)
				}
			}
		})
	}
}

func TestLexer_ComplexProgram(t *testing.T) {
	input := `
// Set seed
seed(42);

// Variable declarations
int(1, 100) x;
float pi = 3.14;
string name = "WTFScript";
bool flag = true;

// Arithmetic
int result = (x + 10) * 2 - 5;

// Print
print(x, pi, name, flag);
`

	expectedTokens := []TokenType{
		// seed(42);
		IDENT, LPAREN, INT, RPAREN, SEMICOLON,

		// int(1, 100) x;
		TYPE_INT, LPAREN, INT, COMMA, INT, RPAREN, IDENT, SEMICOLON,

		// float pi = 3.14;
		TYPE_FLOAT, IDENT, ASSIGN, FLOAT, SEMICOLON,

		// string name = "WTFScript";
		TYPE_STRING, IDENT, ASSIGN, STRING, SEMICOLON,

		// bool flag = true;
		TYPE_BOOL, IDENT, ASSIGN, TRUE, SEMICOLON,

		// int result = (x + 10) * 2 - 5;
		TYPE_INT, IDENT, ASSIGN, LPAREN, IDENT, PLUS, INT, RPAREN, ASTERISK, INT, MINUS, INT, SEMICOLON,

		// print(x, pi, name, flag);
		IDENT, LPAREN, IDENT, COMMA, IDENT, COMMA, IDENT, COMMA, IDENT, RPAREN, SEMICOLON,

		EOF,
	}

	lexer := NewLexer("complex_program", input)

	for i, expectedType := range expectedTokens {
		tok := lexer.NextToken()
		if tok.Type != expectedType {
			t.Errorf("token[%d] - wrong type. expected=%q, got=%q (literal: %q)",
				i, expectedType, tok.Type, tok.Literal)
		}
	}
}

func TestLexer_LineAndColumnTracking(t *testing.T) {
	input := `int x = 5;
float y = 3.14;`

	lexer := NewLexer("line_tracking", input)

	// First line tokens
	tok := lexer.NextToken() // int
	if tok.Line != 1 || tok.Column != 1 {
		t.Errorf("expected line 1, column 1, got line %d, column %d", tok.Line, tok.Column)
	}

	// Skip to second line
	for tok.Type != TYPE_FLOAT {
		tok = lexer.NextToken()
	}

	if tok.Line != 2 {
		t.Errorf("expected line 2, got line %d", tok.Line)
	}
}

func TestLexer_ErrorHandling(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "unterminated string",
			input: `string s = "hello`,
		},
		{
			name:  "invalid character",
			input: "int x = @;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.name, tt.input)

			foundError := false
			for {
				tok := lexer.NextToken()
				if tok.Type == ILLEGAL {
					foundError = true
					break
				}
				if tok.Type == EOF {
					break
				}
			}

			if !foundError {
				t.Errorf("expected to find ILLEGAL token for invalid input")
			}
		})
	}
}

// Benchmark tests
func BenchmarkLexer(b *testing.B) {
	input := `
seed(42);
int(1, 100) x;
float pi = 3.14;
string name = "WTFScript";
bool flag = true;
int result = (x + 10) * 2 - 5;
print(x, pi, name, flag);
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer := NewLexer("benchmark", input)
		for {
			tok := lexer.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}
