package interpreter

import (
	"testing"
	"wtf-script/types"
)

func TestInterpreter_VarDeclaration(t *testing.T) {
	input := "int x = 10;"
	i := NewInterpreter(nil)
	i.Execute(input)

	v, ok := i.GetVariable("x")
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

	v, ok := i.GetVariable("x")
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

	v, ok := i.GetVariable("res")
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

	v, ok := i.GetVariable("x")
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

	v, ok := i.GetVariable("x")
	if !ok {
		t.Fatalf("variable 'x' not found")
	}
	if val, ok := v.Value.(int64); !ok || val != -5 {
		t.Errorf("expected -5, got %v", v.Value)
	}
}
