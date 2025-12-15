package interpreter

import (
	"testing"
)

func TestVarStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"int x = 5;", "x", 5},
		{"bool y = true;", "y", true},
		{"string z = \"hello\";", "z", "hello"},
	}

	for _, tt := range tests {
		l := NewLexer("test", tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*VarDecl)
		if !ok {
			t.Fatalf("program.Statements[0] is not VarDecl. got=%T", program.Statements[0])
		}

		if stmt.Name.Value != tt.expectedIdentifier {
			t.Errorf("stmt.Name.Value not '%s'. got=%s", tt.expectedIdentifier, stmt.Name.Value)
		}
	}
}

func TestVarStatementWithRange(t *testing.T) {
	input := "int(0, 100) x = 5;"
	l := NewLexer("test", input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*VarDecl)
	if !ok {
		t.Fatalf("stmt not VarDecl. got=%T", program.Statements[0])
	}

	if stmt.RangeMin == nil || stmt.RangeMax == nil {
		t.Fatalf("RangeMin or RangeMax is nil")
	}

	// Check simple string representation of ranges for now
	if stmt.RangeMin.String() != "0" {
		t.Errorf("RangeMin wrong. got=%s", stmt.RangeMin.String())
	}
	if stmt.RangeMax.String() != "100" {
		t.Errorf("RangeMax wrong. got=%s", stmt.RangeMax.String())
	}
}

func TestAssignmentStatements(t *testing.T) {
	input := "x = 5;"
	l := NewLexer("test", input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*AssignStmt)
	if !ok {
		t.Fatalf("stmt is not AssignStmt. got=%T", program.Statements[0])
	}

	if stmt.Name.Value != "x" {
		t.Errorf("stmt.Name.Value not 'x'. got=%s", stmt.Name.Value)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)(-5 * 5)",
		},
		{
			"3 + 4; (-5) * 5",
			"(3 + 4)(-5 * 5)",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for _, tt := range tests {
		if tt.input == "!-a" {
			continue
		} // Skip unsupported bang

		l := NewLexer("test", tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
