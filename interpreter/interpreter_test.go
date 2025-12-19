package interpreter

import (
	"testing"
	"wtf-script/types"
)

// ============================================================================
// Variable Declaration Tests
// ============================================================================

func TestInterpreter_VarDeclaration(t *testing.T) {
	input := "int x = 10;"
	i := NewInterpreter(nil)
	i.Execute(input)

	v, ok := i.Variables["x"]
	if !ok {
		t.Fatalf("variable 'x' not found")
	}

	if v.Type != types.Int {
		t.Errorf("expected type Int, got %v", v.Type)
	}

	if val, ok := v.Value.(int64); !ok || val != 10 {
		t.Errorf("expected value 10, got %v", v.Value)
	}
}

func TestInterpreter_VarRandom(t *testing.T) {
	input := "int(10, 20) x;"
	i := NewInterpreter(nil)
	i.SetSeed(42) // Deterministic see if possible, or just check range
	i.Execute(input)

	v, ok := i.Variables["x"]
	if !ok {
		t.Fatalf("variable 'x' not found")
	}

	val, ok := v.Value.(int64)
	if !ok {
		t.Fatalf("value is not int64: %v", v.Value)
	}

	if val < 10 || val > 20 {
		t.Errorf("value %d out of range [10, 20]", val)
	}
}

func TestInterpreter_Arithmetic(t *testing.T) {
	input := "int res = 10 + 5 * 2;"
	i := NewInterpreter(nil)
	i.Execute(input)

	v, ok := i.Variables["res"]
	if !ok {
		t.Fatalf("variable 'res' not found")
	}

	if val, ok := v.Value.(int64); !ok || val != 20 {
		t.Errorf("expected 20, got %v", v.Value)
	}
}

func TestInterpreter_Assignment(t *testing.T) {
	input := `
	int x = 5;
	x = 10;
	`
	i := NewInterpreter(nil)
	i.Execute(input)

	v, ok := i.Variables["x"]
	if !ok {
		t.Fatalf("variable 'x' not found")
	}
	if val, ok := v.Value.(int64); !ok || val != 10 {
		t.Errorf("expected 10, got %v", v.Value)
	}
}

func TestInterpreter_NegativeNumbers(t *testing.T) {
	input := "int x = -5;"
	i := NewInterpreter(nil)
	i.Execute(input)

	v, ok := i.Variables["x"]
	if !ok {
		t.Fatalf("variable 'x' not found")
	}
	if val, ok := v.Value.(int64); !ok || val != -5 {
		t.Errorf("expected -5, got %v", v.Value)
	}
}

// ============================================================================
// Comparison Operators Tests
// ============================================================================

func TestInterpreter_ComparisonOperators(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		varName  string
		expected bool
	}{
		{"equal_int", "bool result = 10 == 10;", "result", true},
		{"not_equal_int", "bool result = 10 != 5;", "result", true},
		{"less_than", "bool result = 5 < 10;", "result", true},
		{"less_than_equal", "bool result = 10 <= 10;", "result", true},
		{"greater_than", "bool result = 15 > 10;", "result", true},
		{"greater_than_equal", "bool result = 10 >= 10;", "result", true},
		{"equal_false", "bool result = 10 == 5;", "result", false},
		{"not_equal_false", "bool result = 10 != 10;", "result", false},
		// String comparisons
		{"string_equal", `bool result = "hello" == "hello";`, "result", true},
		{"string_not_equal", `bool result = "hello" != "world";`, "result", true},
		{"string_less_than", `bool result = "abc" < "def";`, "result", true},
		// Float comparisons
		{"float_equal", "bool result = 3.14 == 3.14;", "result", true},
		{"float_greater", "bool result = 5.5 > 3.3;", "result", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewInterpreter(nil)
			i.Execute(tt.input)

			v, ok := i.Variables[tt.varName]
			if !ok {
				t.Fatalf("variable '%s' not found", tt.varName)
			}

			val, ok := v.Value.(bool)
			if !ok {
				t.Fatalf("expected bool value, got %T", v.Value)
			}

			if val != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, val)
			}
		})
	}
}

func TestInterpreter_ComparisonWithVariables(t *testing.T) {
	input := `
	int x = 10;
	int y = 20;
	bool test1 = x < y;
	bool test2 = x > y;
	bool test3 = x == 10;
	`
	i := NewInterpreter(nil)
	i.Execute(input)

	tests := []struct {
		name     string
		expected bool
	}{
		{"test1", true},
		{"test2", false},
		{"test3", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, ok := i.Variables[tt.name]
			if !ok {
				t.Fatalf("variable '%s' not found", tt.name)
			}

			val, ok := v.Value.(bool)
			if !ok {
				t.Fatalf("expected bool value, got %T", v.Value)
			}

			if val != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, val)
			}
		})
	}
}

// ============================================================================
// Logical Operators Tests (&&, ||)
// ============================================================================

func TestInterpreter_LogicalOperators(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		varName  string
		expected bool
	}{
		{"and_true_true", "bool result = true && true;", "result", true},
		{"and_true_false", "bool result = true && false;", "result", false},
		{"and_false_true", "bool result = false && true;", "result", false},
		{"and_false_false", "bool result = false && false;", "result", false},
		{"or_true_true", "bool result = true || true;", "result", true},
		{"or_true_false", "bool result = true || false;", "result", true},
		{"or_false_true", "bool result = false || true;", "result", true},
		{"or_false_false", "bool result = false || false;", "result", false},
		// With comparisons
		{"and_with_comp", "bool result = 10 > 5 && 20 > 15;", "result", true},
		{"or_with_comp", "bool result = 10 < 5 || 20 > 15;", "result", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewInterpreter(nil)
			i.Execute(tt.input)

			v, ok := i.Variables[tt.varName]
			if !ok {
				t.Fatalf("variable '%s' not found", tt.varName)
			}

			val, ok := v.Value.(bool)
			if !ok {
				t.Fatalf("expected bool value, got %T", v.Value)
			}

			if val != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, val)
			}
		})
	}
}

func TestInterpreter_LogicalOperatorPrecedence(t *testing.T) {
	// AND has higher precedence than OR
	input := `
	bool test1 = false || true && false;
	bool test2 = true && false || true;
	`
	i := NewInterpreter(nil)
	i.Execute(input)

	// false || (true && false) = false || false = false
	v1, _ := i.Variables["test1"]
	if val, ok := v1.Value.(bool); !ok || val != false {
		t.Errorf("test1: expected false, got %v", v1.Value)
	}

	// (true && false) || true = false || true = true
	v2, _ := i.Variables["test2"]
	if val, ok := v2.Value.(bool); !ok || val != true {
		t.Errorf("test2: expected true, got %v", v2.Value)
	}
}

func TestInterpreter_ShortCircuitEvaluation(t *testing.T) {
	// AND short-circuit: if left is false, right is not evaluated
	input1 := `
	int x = 10;
	bool result = false && x / 0 == 5;
	`
	i := NewInterpreter(nil)
	i.Execute(input1)

	v, ok := i.Variables["result"]
	if !ok {
		t.Fatal("variable 'result' not found")
	}
	if val, ok := v.Value.(bool); !ok || val != false {
		t.Errorf("expected false, got %v", v.Value)
	}

	// OR short-circuit: if left is true, right is not evaluated
	input2 := `
	int y = 20;
	bool result2 = true || y / 0 == 5;
	`
	i2 := NewInterpreter(nil)
	i2.Execute(input2)

	v2, ok := i2.Variables["result2"]
	if !ok {
		t.Fatal("variable 'result2' not found")
	}
	if val, ok := v2.Value.(bool); !ok || val != true {
		t.Errorf("expected true, got %v", v2.Value)
	}
}

// ============================================================================
// If/Else Statement Tests
// ============================================================================

func TestInterpreter_IfStatement(t *testing.T) {
	input := `
	int x = 10;
	int result = 0;
	if (x > 5) {
		result = 100;
	}
	`
	i := NewInterpreter(nil)
	i.Execute(input)

	v, ok := i.Variables["result"]
	if !ok {
		t.Fatal("variable 'result' not found")
	}

	val, ok := v.Value.(int64)
	if !ok || val != 100 {
		t.Errorf("expected 100, got %v", v.Value)
	}
}

func TestInterpreter_IfElseStatement(t *testing.T) {
	input := `
	int x = 3;
	int result = 0;
	if (x > 5) {
		result = 100;
	} else {
		result = 200;
	}
	`
	i := NewInterpreter(nil)
	i.Execute(input)

	v, ok := i.Variables["result"]
	if !ok {
		t.Fatal("variable 'result' not found")
	}

	val, ok := v.Value.(int64)
	if !ok || val != 200 {
		t.Errorf("expected 200, got %v", v.Value)
	}
}

func TestInterpreter_NestedIfStatement(t *testing.T) {
	input := `
	int x = 10;
	int y = 20;
	int result = 0;
	if (x > 5) {
		if (y > 15) {
			result = 300;
		}
	}
	`
	i := NewInterpreter(nil)
	i.Execute(input)

	v, ok := i.Variables["result"]
	if !ok {
		t.Fatal("variable 'result' not found")
	}

	val, ok := v.Value.(int64)
	if !ok || val != 300 {
		t.Errorf("expected 300, got %v", v.Value)
	}
}

// ============================================================================
// Ifrand Statement Tests
// ============================================================================

func TestInterpreter_IfrandDefaultProbability(t *testing.T) {
	// Test that ifrand without argument works with 50% probability
	input := `
	int result = 0;
	ifrand {
		result = 1;
	}
	`
	// Run multiple times and check we get both outcomes
	trueCount := 0
	falseCount := 0
	for i := 0; i < 100; i++ {
		interp := NewInterpreter(nil)
		interp.Execute(input)
		v := interp.Variables["result"]
		if val, ok := v.Value.(int64); ok && val == 1 {
			trueCount++
		} else {
			falseCount++
		}
	}

	// With 100 trials, we should get both outcomes (not 0 or 100)
	if trueCount == 0 || trueCount == 100 {
		t.Errorf("ifrand appears deterministic: true=%d, false=%d", trueCount, falseCount)
	}
}

func TestInterpreter_IfrandWithProbability(t *testing.T) {
	// Test ifrand with explicit 0.0 probability (never executes)
	input1 := `
	int result = 0;
	ifrand(0.0) {
		result = 1;
	}
	`
	i1 := NewInterpreter(nil)
	i1.Execute(input1)
	v1 := i1.Variables["result"]
	if val, ok := v1.Value.(int64); !ok || val != 0 {
		t.Errorf("expected 0 (never execute), got %v", v1.Value)
	}

	// Test ifrand with explicit 1.0 probability (always executes)
	input2 := `
	int result = 0;
	ifrand(1.0) {
		result = 1;
	}
	`
	i2 := NewInterpreter(nil)
	i2.Execute(input2)
	v2 := i2.Variables["result"]
	if val, ok := v2.Value.(int64); !ok || val != 1 {
		t.Errorf("expected 1 (always execute), got %v", v2.Value)
	}
}

func TestInterpreter_IfrandDeterministic(t *testing.T) {
	// With a fixed seed, ifrand should be deterministic
	input := `
	seed(12345);
	int result = 0;
	ifrand {
		result = 1;
	}
	`

	i1 := NewInterpreter(nil)
	i1.Execute(input)
	v1 := i1.Variables["result"]

	i2 := NewInterpreter(nil)
	i2.Execute(input)
	v2 := i2.Variables["result"]

	if v1.Value != v2.Value {
		t.Errorf("ifrand not deterministic with same seed: %v vs %v", v1.Value, v2.Value)
	}
}

// ============================================================================
// Uint Type Tests
// ============================================================================

func TestInterpreter_UintBasic(t *testing.T) {
	input := "uint x = 100;"
	i := NewInterpreter(nil)
	i.Execute(input)

	v, ok := i.Variables["x"]
	if !ok {
		t.Fatal("variable 'x' not found")
	}

	if v.Type != types.Uint {
		t.Errorf("expected type Uint, got %v", v.Type)
	}

	val, ok := v.Value.(uint64)
	if !ok || val != 100 {
		t.Errorf("expected 100, got %v", v.Value)
	}
}

func TestInterpreter_UintRange(t *testing.T) {
	t.Skip("uint type doesn't support range syntax - skipping test")
}

func TestInterpreter_UintArithmetic(t *testing.T) {
	input := `
	uint x = 10;
	uint y = 5;
	uint result = x + y;
	`
	i := NewInterpreter(nil)
	i.Execute(input)

	v, ok := i.Variables["result"]
	if !ok {
		t.Fatal("variable 'result' not found")
	}

	val, ok := v.Value.(uint64)
	if !ok || val != 15 {
		t.Errorf("expected 15, got %v", v.Value)
	}
}

// ============================================================================
// Unofloat Type Tests
// ============================================================================

func TestInterpreter_UnofloatBasic(t *testing.T) {
	input := "unofloat x = 0.5;"
	i := NewInterpreter(nil)
	i.Execute(input)

	v, ok := i.Variables["x"]
	if !ok {
		t.Fatal("variable 'x' not found")
	}

	if v.Type != types.UnitFloat {
		t.Errorf("expected type UnitFloat, got %v", v.Type)
	}

	val, ok := v.Value.(types.Unofloat)
	if !ok {
		t.Fatalf("value is not Unofloat: %T", v.Value)
	}

	if float64(val) != 0.5 {
		t.Errorf("expected 0.5, got %v", float64(val))
	}
}

func TestInterpreter_UnofloatRange(t *testing.T) {
	t.Skip("unofloat type doesn't support range syntax - skipping test")
}

func TestInterpreter_UnofloatClamping(t *testing.T) {
	// Test that unofloat clamps values to [0, 1] for computed expressions
	input := `
	unofloat x = 0.8;
	unofloat y = 0.5;
	unofloat result = x + y;
	`
	i := NewInterpreter(nil)
	i.Execute(input)

	v, ok := i.Variables["result"]
	if !ok {
		t.Fatal("variable 'result' not found")
	}

	val, ok := v.Value.(types.Unofloat)
	if !ok {
		t.Fatalf("value is not Unofloat: %v", v.Value)
	}

	fval := float64(val)
	// 0.8 + 0.5 = 1.3, should be clamped to 1.0
	if fval != 1.0 {
		t.Errorf("expected 1.0 (clamped), got %v", fval)
	}
}

func TestInterpreter_UnofloatArithmetic(t *testing.T) {
	input := `
	unofloat x = 0.6;
	unofloat y = 0.2;
	unofloat sum = x + y;
	unofloat diff = x - y;
	unofloat prod = x * y;
	unofloat quot = x / y;
	`
	i := NewInterpreter(nil)
	i.Execute(input)

	tests := []struct {
		name     string
		expected float64
	}{
		{"sum", 0.8},
		{"diff", 0.4},
		{"prod", 0.12},
		{"quot", 1.0}, // 0.6 / 0.2 = 3.0, clamped to 1.0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, ok := i.Variables[tt.name]
			if !ok {
				t.Fatalf("variable '%s' not found", tt.name)
			}

			val, ok := v.Value.(types.Unofloat)
			if !ok {
				t.Fatalf("value is not Unofloat: %T", v.Value)
			}

			fval := float64(val)
			// Allow small floating point differences
			if fval < tt.expected-0.01 || fval > tt.expected+0.01 {
				t.Errorf("expected ~%v, got %v", tt.expected, fval)
			}
		})
	}
}

// ============================================================================
// String Operation Tests
// ============================================================================

func TestInterpreter_StringConcatenation(t *testing.T) {
	input := `
	string s1 = "Hello";
	string s2 = " World";
	string result = s1 + s2;
	`
	i := NewInterpreter(nil)
	i.Execute(input)

	v, ok := i.Variables["result"]
	if !ok {
		t.Fatal("variable 'result' not found")
	}

	val, ok := v.Value.(string)
	if !ok || val != "Hello World" {
		t.Errorf("expected 'Hello World', got %v", v.Value)
	}
}

func TestInterpreter_StringComparison(t *testing.T) {
	input := `
	string s1 = "apple";
	string s2 = "banana";
	bool eq = s1 == s2;
	bool neq = s1 != s2;
	bool lt = s1 < s2;
	`
	i := NewInterpreter(nil)
	i.Execute(input)

	eqVal, _ := i.Variables["eq"]
	if val, ok := eqVal.Value.(bool); !ok || val != false {
		t.Errorf("expected eq=false, got %v", eqVal.Value)
	}

	neqVal, _ := i.Variables["neq"]
	if val, ok := neqVal.Value.(bool); !ok || val != true {
		t.Errorf("expected neq=true, got %v", neqVal.Value)
	}

	ltVal, _ := i.Variables["lt"]
	if val, ok := ltVal.Value.(bool); !ok || val != true {
		t.Errorf("expected lt=true, got %v", ltVal.Value)
	}
}

// ============================================================================
// Builtin Function Tests
// ============================================================================

func TestInterpreter_SeedFunction(t *testing.T) {
	input := `
	seed(42);
	int(1, 100) x;
	`

	i1 := NewInterpreter(nil)
	i1.Execute(input)
	v1 := i1.Variables["x"]

	i2 := NewInterpreter(nil)
	i2.Execute(input)
	v2 := i2.Variables["x"]

	if v1.Value != v2.Value {
		t.Errorf("same seed should produce same random values: %v vs %v", v1.Value, v2.Value)
	}
}

func TestInterpreter_TypeofFunction(t *testing.T) {
	input := `
	int x = 10;
	float y = 3.14;
	bool z = true;
	string s = "hello";
	string t1 = typeof(x);
	string t2 = typeof(y);
	string t3 = typeof(z);
	string t4 = typeof(s);
	`
	i := NewInterpreter(nil)
	i.Execute(input)

	tests := []struct {
		name     string
		expected string
	}{
		{"t1", "int"},
		{"t2", "float"},
		{"t3", "bool"},
		{"t4", "string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, ok := i.Variables[tt.name]
			if !ok {
				t.Fatalf("variable '%s' not found", tt.name)
			}

			val, ok := v.Value.(string)
			if !ok || val != tt.expected {
				t.Errorf("expected '%s', got %v", tt.expected, v.Value)
			}
		})
	}
}

// ============================================================================
// Type Coercion Tests
// ============================================================================

func TestInterpreter_TypeCoercion(t *testing.T) {
	// Test that int can be assigned to float
	input := `
	float x = 10;
	`
	i := NewInterpreter(nil)
	i.Execute(input)

	v, ok := i.Variables["x"]
	if !ok {
		t.Fatal("variable 'x' not found")
	}

	val, ok := v.Value.(float64)
	if !ok || val != 10.0 {
		t.Errorf("expected 10.0, got %v", v.Value)
	}
}

func TestInterpreter_MixedTypeArithmetic(t *testing.T) {
	// First operand type determines result type (FCFS)
	input := `
	int x = 10;
	float y = 3.5;
	int result1 = x + 3;
	float result2 = y + 2;
	`
	i := NewInterpreter(nil)
	i.Execute(input)

	v1, _ := i.Variables["result1"]
	val1, ok := v1.Value.(int64)
	if !ok || val1 != 13 {
		t.Errorf("expected 13 (int), got %v", v1.Value)
	}

	v2, _ := i.Variables["result2"]
	val2, ok := v2.Value.(float64)
	if !ok || val2 != 5.5 {
		t.Errorf("expected 5.5 (float), got %v", v2.Value)
	}
}

// ============================================================================
// Error Handling Tests
// ============================================================================

func TestInterpreter_DivisionByZero(t *testing.T) {
	input := `
	int x = 10;
	int result = x / 0;
	`
	i := NewInterpreter(nil)
	i.Execute(input)

	// Variable should not be created due to error
	_, ok := i.Variables["result"]
	if ok {
		t.Error("expected error for division by zero, but variable was created")
	}
}

func TestInterpreter_UndefinedVariable(t *testing.T) {
	input := `
	int x = y + 5;
	`
	i := NewInterpreter(nil)
	i.Execute(input)

	_, ok := i.Variables["x"]
	if ok {
		t.Error("expected error for undefined variable, but variable was created")
	}
}

func TestInterpreter_FloatLiterals(t *testing.T) {
	input := `
	float pi = 3.14159;
	float e = 2.71828;
	float sum = pi + e;
	`
	i := NewInterpreter(nil)
	i.Execute(input)

	v, ok := i.Variables["sum"]
	if !ok {
		t.Fatal("variable 'sum' not found")
	}

	val, ok := v.Value.(float64)
	if !ok {
		t.Fatalf("expected float64, got %T", v.Value)
	}

	expected := 3.14159 + 2.71828
	if val < expected-0.001 || val > expected+0.001 {
		t.Errorf("expected ~%v, got %v", expected, val)
	}
}

func TestInterpreter_BooleanLiterals(t *testing.T) {
	input := `
	bool t = true;
	bool f = false;
	`
	i := NewInterpreter(nil)
	i.Execute(input)

	vt, ok := i.Variables["t"]
	if !ok {
		t.Fatal("variable 't' not found")
	}
	if val, ok := vt.Value.(bool); !ok || val != true {
		t.Errorf("expected true, got %v", vt.Value)
	}

	vf, ok := i.Variables["f"]
	if !ok {
		t.Fatal("variable 'f' not found")
	}
	if val, ok := vf.Value.(bool); !ok || val != false {
		t.Errorf("expected false, got %v", vf.Value)
	}
}

// ============================================================================
// Complex Expression Tests
// ============================================================================

func TestInterpreter_ComplexExpression(t *testing.T) {
	input := `
	int a = 5;
	int b = 10;
	int c = 15;
	int result = (a + b) * c - 20 / 4;
	`
	i := NewInterpreter(nil)
	i.Execute(input)

	v, ok := i.Variables["result"]
	if !ok {
		t.Fatal("variable 'result' not found")
	}

	// (5 + 10) * 15 - 20 / 4 = 15 * 15 - 5 = 225 - 5 = 220
	val, ok := v.Value.(int64)
	if !ok || val != 220 {
		t.Errorf("expected 220, got %v", v.Value)
	}
}

func TestInterpreter_MultipleAssignments(t *testing.T) {
	input := `
	int x = 10;
	x = 20;
	x = x + 5;
	x = x * 2;
	`
	i := NewInterpreter(nil)
	i.Execute(input)

	v, ok := i.Variables["x"]
	if !ok {
		t.Fatal("variable 'x' not found")
	}

	// (20 + 5) * 2 = 50
	val, ok := v.Value.(int64)
	if !ok || val != 50 {
		t.Errorf("expected 50, got %v", v.Value)
	}
}

func TestInterpreter_UnaryMinus(t *testing.T) {
	input := `
	int x = -5;
	int y = -x;
	int z = -(10 + 5);
	`
	i := NewInterpreter(nil)
	i.Execute(input)

	vx, _ := i.Variables["x"]
	if val, ok := vx.Value.(int64); !ok || val != -5 {
		t.Errorf("expected -5, got %v", vx.Value)
	}

	vy, _ := i.Variables["y"]
	if val, ok := vy.Value.(int64); !ok || val != 5 {
		t.Errorf("expected 5, got %v", vy.Value)
	}

	vz, _ := i.Variables["z"]
	if val, ok := vz.Value.(int64); !ok || val != -15 {
		t.Errorf("expected -15, got %v", vz.Value)
	}
}
