package interpreter

import (
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&VarDecl{
				Token: Token{Type: TYPE_INT, Literal: "int"},
				Name: &Identifier{
					Token: Token{Type: IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: Token{Type: IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "int myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

func TestStringWithRange(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&VarDecl{
				Token: Token{Type: TYPE_INT, Literal: "int"},
				Name: &Identifier{
					Token: Token{Type: IDENT, Literal: "x"},
					Value: "x",
				},
				RangeMin: &IntegerLiteral{
					Token: Token{Type: INT, Literal: "0"},
					Value: 0,
				},
				RangeMax: &IntegerLiteral{
					Token: Token{Type: INT, Literal: "100"},
					Value: 100,
				},
				Value: &IntegerLiteral{
					Token: Token{Type: INT, Literal: "5"},
					Value: 5,
				},
			},
		},
	}

	if program.String() != "int(0, 100) x = 5;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
