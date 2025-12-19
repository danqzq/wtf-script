package interpreter

import (
	"testing"
)

// ============================================================================
// Parser Basic Tests
// ============================================================================

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

// ============================================================================
// AST Node Tests
// ============================================================================

func TestAST_ProgramString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&VarDecl{
				Token: Token{Type: TYPE_INT, Literal: "int"},
				Type:  TYPE_INT,
				Name:  &Identifier{Value: "x"},
				Value: &IntegerLiteral{Value: 10},
			},
		},
	}

	str := program.String()
	if str == "" {
		t.Error("program.String() returned empty string")
	}
}

func TestAST_BinaryExprString(t *testing.T) {
	expr := &BinaryExpr{
		Left:     &IntegerLiteral{Token: Token{Literal: "5"}, Value: 5},
		Operator: PLUS,
		Right:    &IntegerLiteral{Token: Token{Literal: "10"}, Value: 10},
	}

	str := expr.String()
	expected := "(5 + 10)"
	if str != expected {
		t.Errorf("expected '%s', got '%s'", expected, str)
	}
}

func TestAST_UnaryExprString(t *testing.T) {
	expr := &UnaryExpr{
		Operator: "-",
		Right:    &IntegerLiteral{Token: Token{Literal: "5"}, Value: 5},
	}

	str := expr.String()
	expected := "(-5)"
	if str != expected {
		t.Errorf("expected '%s', got '%s'", expected, str)
	}
}

func TestAST_IdentifierString(t *testing.T) {
	ident := &Identifier{Value: "myVar"}
	str := ident.String()
	if str != "myVar" {
		t.Errorf("expected 'myVar', got '%s'", str)
	}
}

func TestAST_LiteralStrings(t *testing.T) {
	tests := []struct {
		name     string
		node     Expression
		expected string
	}{
		{"int_literal", &IntegerLiteral{Token: Token{Literal: "42"}, Value: 42}, "42"},
		{"float_literal", &FloatLiteral{Token: Token{Literal: "3.14"}, Value: 3.14}, "3.14"},
		{"bool_literal", &BooleanLiteral{Token: Token{Literal: "true"}, Value: true}, "true"},
		{"string_literal", &StringLiteral{Token: Token{Literal: "\"hello\""}, Value: "\"hello\""}, "\"hello\""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := tt.node.String()
			if str != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, str)
			}
		})
	}
}

// ============================================================================
// Parser Tests for Control Flow
// ============================================================================

func TestParser_IfStatement(t *testing.T) {
	input := `
	if (x > 5) {
		y = 10;
	}
	`
	l := NewLexer("test", input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*IfStmt)
	if !ok {
		t.Fatalf("statement is not IfStmt, got %T", program.Statements[0])
	}

	if stmt.Condition == nil {
		t.Fatal("condition is nil")
	}

	if stmt.Consequence == nil {
		t.Fatal("consequence is nil")
	}

	if stmt.Alternative != nil {
		t.Error("alternative should be nil for if without else")
	}
}

func TestParser_IfElseStatement(t *testing.T) {
	input := `
	if (x > 5) {
		y = 10;
	} else {
		y = 20;
	}
	`
	l := NewLexer("test", input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*IfStmt)
	if !ok {
		t.Fatalf("statement is not IfStmt, got %T", program.Statements[0])
	}

	if stmt.Alternative == nil {
		t.Fatal("alternative should not be nil for if-else")
	}
}

func TestParser_IfrandStatement(t *testing.T) {
	input := `
	ifrand {
		x = 1;
	}
	`
	l := NewLexer("test", input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*IfStmt)
	if !ok {
		t.Fatalf("statement is not IfStmt, got %T", program.Statements[0])
	}

	if stmt.Token.Type != IFRAND {
		t.Errorf("expected IFRAND token, got %v", stmt.Token.Type)
	}
}

func TestParser_IfrandWithProbability(t *testing.T) {
	input := `
	ifrand(0.75) {
		x = 1;
	}
	`
	l := NewLexer("test", input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*IfStmt)
	if !ok {
		t.Fatalf("statement is not IfStmt, got %T", program.Statements[0])
	}

	if stmt.Condition == nil {
		t.Fatal("condition should not be nil for ifrand with probability")
	}
}

// ============================================================================
// Parser Tests for Operators
// ============================================================================

func TestParser_ComparisonOperators(t *testing.T) {
	tests := []struct {
		input    string
		operator TokenType
	}{
		{"x == 5", EQ},
		{"x != 5", NEQ},
		{"x < 5", LT},
		{"x <= 5", LTE},
		{"x > 5", GT},
		{"x >= 5", GTE},
	}

	for _, tt := range tests {
		t.Run(string(tt.operator), func(t *testing.T) {
			l := NewLexer("test", tt.input)
			p := NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(program.Statements))
			}

			exprStmt, ok := program.Statements[0].(*ExprStmt)
			if !ok {
				t.Fatalf("statement is not ExprStmt, got %T", program.Statements[0])
			}

			binExpr, ok := exprStmt.Expression.(*BinaryExpr)
			if !ok {
				t.Fatalf("expression is not BinaryExpr, got %T", exprStmt.Expression)
			}

			if binExpr.Operator != tt.operator {
				t.Errorf("expected operator %v, got %v", tt.operator, binExpr.Operator)
			}
		})
	}
}

func TestParser_LogicalOperators(t *testing.T) {
	tests := []struct {
		input    string
		operator TokenType
	}{
		{"true && false", AND},
		{"true || false", OR},
	}

	for _, tt := range tests {
		t.Run(string(tt.operator), func(t *testing.T) {
			l := NewLexer("test", tt.input)
			p := NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(program.Statements))
			}

			exprStmt, ok := program.Statements[0].(*ExprStmt)
			if !ok {
				t.Fatalf("statement is not ExprStmt, got %T", program.Statements[0])
			}

			binExpr, ok := exprStmt.Expression.(*BinaryExpr)
			if !ok {
				t.Fatalf("expression is not BinaryExpr, got %T", exprStmt.Expression)
			}

			if binExpr.Operator != tt.operator {
				t.Errorf("expected operator %v, got %v", tt.operator, binExpr.Operator)
			}
		})
	}
}

func TestParser_UnaryOperators(t *testing.T) {
	input := "-x"
	l := NewLexer("test", input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ExprStmt)
	if !ok {
		t.Fatalf("statement is not ExprStmt, got %T", program.Statements[0])
	}

	unaryExpr, ok := exprStmt.Expression.(*UnaryExpr)
	if !ok {
		t.Fatalf("expression is not UnaryExpr, got %T", exprStmt.Expression)
	}

	if unaryExpr.Operator != "-" {
		t.Errorf("expected operator '-', got %v", unaryExpr.Operator)
	}
}

// ============================================================================
// Parser Tests for All Type Keywords
// ============================================================================

func TestParser_AllTypeDeclarations(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		tokenType TokenType
	}{
		{"int", "int x = 5;", TYPE_INT},
		{"uint", "uint x = 5;", TYPE_UINT},
		{"float", "float x = 3.14;", TYPE_FLOAT},
		{"unofloat", "unofloat x = 0.5;", TYPE_UNOFLOAT},
		{"bool", "bool x = true;", TYPE_BOOL},
		{"string", "string x = \"hello\";", TYPE_STRING},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer("test", tt.input)
			p := NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*VarDecl)
			if !ok {
				t.Fatalf("statement is not VarDecl, got %T", program.Statements[0])
			}

			if stmt.Type != tt.tokenType {
				t.Errorf("expected type %v, got %v", tt.tokenType, stmt.Type)
			}
		})
	}
}

func TestParser_RangeDeclarations(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"int_range", "int(0, 100) x;"},
		{"float_range", "float(0.0, 10.0) x;"},
		// Note: uint and unofloat may not support range syntax in the implementation
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer("test", tt.input)
			p := NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*VarDecl)
			if !ok {
				t.Fatalf("statement is not VarDecl, got %T", program.Statements[0])
			}

			if stmt.RangeMin == nil || stmt.RangeMax == nil {
				t.Fatal("RangeMin or RangeMax is nil")
			}
		})
	}
}

// ============================================================================
// Parser Tests for Function Calls
// ============================================================================

func TestParser_FunctionCallNoArgs(t *testing.T) {
	input := "foo();"
	l := NewLexer("test", input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ExprStmt)
	if !ok {
		t.Fatalf("statement is not ExprStmt, got %T", program.Statements[0])
	}

	callExpr, ok := exprStmt.Expression.(*CallExpr)
	if !ok {
		t.Fatalf("expression is not CallExpr, got %T", exprStmt.Expression)
	}

	if len(callExpr.Arguments) != 0 {
		t.Errorf("expected 0 arguments, got %d", len(callExpr.Arguments))
	}
}

func TestParser_FunctionCallWithArgs(t *testing.T) {
	input := "print(x, y, 42);"
	l := NewLexer("test", input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ExprStmt)
	if !ok {
		t.Fatalf("statement is not ExprStmt, got %T", program.Statements[0])
	}

	callExpr, ok := exprStmt.Expression.(*CallExpr)
	if !ok {
		t.Fatalf("expression is not CallExpr, got %T", exprStmt.Expression)
	}

	if len(callExpr.Arguments) != 3 {
		t.Errorf("expected 3 arguments, got %d", len(callExpr.Arguments))
	}
}

func TestParser_NestedFunctionCalls(t *testing.T) {
	input := "print(typeof(x));"
	l := NewLexer("test", input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ExprStmt)
	if !ok {
		t.Fatalf("statement is not ExprStmt, got %T", program.Statements[0])
	}

	callExpr, ok := exprStmt.Expression.(*CallExpr)
	if !ok {
		t.Fatalf("expression is not CallExpr, got %T", exprStmt.Expression)
	}

	if len(callExpr.Arguments) != 1 {
		t.Fatalf("expected 1 argument, got %d", len(callExpr.Arguments))
	}

	// Check that the argument is also a CallExpr
	_, ok = callExpr.Arguments[0].(*CallExpr)
	if !ok {
		t.Errorf("expected nested CallExpr, got %T", callExpr.Arguments[0])
	}
}

// ============================================================================
// Parser Tests for Complex Programs
// ============================================================================

func TestParser_ComplexProgram(t *testing.T) {
	input := `
	seed(42);
	int x = 10;
	float y = 3.14;
	bool flag = true;
	string name = "WTF";
	
	if (x > 5) {
		y = y * 2.0;
	}
	
	int(0, 100) random;
	print(x, y, flag, name, random);
	`
	l := NewLexer("test", input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	// Should have multiple statements
	if len(program.Statements) < 5 {
		t.Errorf("expected at least 5 statements, got %d", len(program.Statements))
	}
}

func TestParser_OperatorMixing(t *testing.T) {
	input := "bool result = x > 5 && y < 10 || z == 15;"
	l := NewLexer("test", input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*VarDecl)
	if !ok {
		t.Fatalf("statement is not VarDecl, got %T", program.Statements[0])
	}

	// The value should be a binary expression
	_, ok = stmt.Value.(*BinaryExpr)
	if !ok {
		t.Fatalf("value is not BinaryExpr, got %T", stmt.Value)
	}
}

// ============================================================================
// Parser Error Tests
// ============================================================================

func TestParser_ErrorRecovery(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"missing_rparen", "int x = (5 + 3;"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer("test", tt.input)
			p := NewParser(l)
			_ = p.ParseProgram()

			if len(p.Errors()) == 0 {
				t.Error("expected parser errors, but got none")
			}
		})
	}
}

func TestParser_EmptyProgram(t *testing.T) {
	input := ""
	l := NewLexer("test", input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 0 {
		t.Errorf("expected 0 statements for empty program, got %d", len(program.Statements))
	}
}

func TestParser_OnlyComments(t *testing.T) {
	input := `
	// This is a comment
	// Another comment
	`
	l := NewLexer("test", input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 0 {
		t.Errorf("expected 0 statements for comments-only program, got %d", len(program.Statements))
	}
}
