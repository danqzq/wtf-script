package interpreter

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"wtf-script/builtins"
	"wtf-script/config"
	"wtf-script/types"
)

type Interpreter struct {
	Variables map[string]types.Variable
	Builtins  map[string]types.IBuiltinFunc
	Rand      *rand.Rand
	Config    *config.Config
}

func (i *Interpreter) GetVariable(name string) (*types.Variable, bool) {
	v, ok := i.Variables[name]
	return &v, ok
}

func (i *Interpreter) SetVariable(name string, value any) {
	if v, ok := i.Variables[name]; ok {
		v.Value = value
		i.Variables[name] = v
	}
}

func (i *Interpreter) SetSeed(seed int64) {
	i.Rand.Seed(seed)
}

func (i *Interpreter) GenerateRandomString(n int, charset string) string {
	b := make([]byte, n)
	for j := range b {
		b[j] = charset[i.Rand.Intn(len(charset))]
	}
	return string(b)
}

func (i *Interpreter) GetConfig() *config.Config {
	return i.Config
}

func NewInterpreter(cfg *config.Config) *Interpreter {
	if cfg == nil {
		cfg = &config.DefaultConfig
	}
	i := &Interpreter{
		Variables: make(map[string]types.Variable),
		Builtins:  make(map[string]types.IBuiltinFunc),
		Rand:      rand.New(rand.NewSource(time.Now().UnixNano())),
		Config:    cfg,
	}

	builtins.RegisterBuiltins(func(name string, fn types.IBuiltinFunc) {
		i.Builtins[name] = fn
	})
	return i
}

func (i *Interpreter) Execute(code string) {
	l := NewLexer("main", code)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, msg := range p.Errors() {
			fmt.Printf("Parser Error: %s\n", msg)
		}
		return
	}

	_, err := i.Evaluate(program)
	if err != nil {
		fmt.Printf("Runtime Error: %s\n", err)
	}
}

func (i *Interpreter) Evaluate(node Node) (any, error) {
	switch node := node.(type) {
	// Program
	case *Program:
		return i.evalProgram(node)

	// Statements
	case *ExprStmt:
		return i.Evaluate(node.Expression)
	case *VarDecl:
		return i.evalVarDecl(node)
	case *AssignStmt:
		return i.evalAssignStmt(node)

	// Expressions
	case *Identifier:
		return i.evalIdentifier(node)
	case *IntegerLiteral:
		return node.Value, nil
	case *FloatLiteral:
		return node.Value, nil
	case *BooleanLiteral:
		return node.Value, nil
	case *StringLiteral:
		if unquoted, err := strconv.Unquote(node.Value); err == nil {
			return unquoted, nil
		}
		return node.Value, nil
	case *BinaryExpr:
		return i.evalBinaryExpr(node)
	case *UnaryExpr:
		return i.evalUnaryExpr(node)
	case *CallExpr:
		return i.evalCallExpr(node)
	}

	return nil, nil
}

func (i *Interpreter) evalProgram(program *Program) (any, error) {
	var result any
	for _, statement := range program.Statements {
		val, err := i.Evaluate(statement)
		if err != nil {
			return nil, err
		}
		result = val
	}
	return result, nil
}

func (i *Interpreter) evalIdentifier(node *Identifier) (any, error) {
	if val, ok := i.Variables[node.Value]; ok {
		return val.Value, nil
	}
	return nil, fmt.Errorf("identifier not found: %s", node.Value)
}

func (i *Interpreter) evalVarDecl(node *VarDecl) (any, error) {
	var val any

	// Handles: int(0, 100) x;
	if node.RangeMin != nil && node.RangeMax != nil {
		minVal, err := i.Evaluate(node.RangeMin)
		if err != nil {
			return nil, err
		}
		maxVal, err := i.Evaluate(node.RangeMax)
		if err != nil {
			return nil, err
		}

		var errVal error
		val, errVal = i.randomValueInRange(node.Type, minVal, maxVal)
		if errVal != nil {
			return nil, errVal
		}
	} else if node.Value != nil {
		// Handles: int x = 5;
		evaluated, err := i.Evaluate(node.Value)
		if err != nil {
			return nil, err
		}
		val = evaluated
	} else {
		// Handles: int x; (random default)
		val = i.randomValue(node.Type)
	}

	// Store variable
	i.Variables[node.Name.Value] = types.Variable{
		Type:  types.VarType(varTypeFromToken(node.Type)), // Helper needed
		Value: val,
	}
	return val, nil // Return value of declaration?
}

func (i *Interpreter) evalAssignStmt(node *AssignStmt) (any, error) {
	val, err := i.Evaluate(node.Value)
	if err != nil {
		return nil, err
	}

	if v, ok := i.Variables[node.Name.Value]; ok {
		// In a real strict language we'd check types here
		// But let's allow implicit casting or just update
		v.Value = val
		i.Variables[node.Name.Value] = v
		return val, nil
	}
	return nil, fmt.Errorf("variable not defined: %s", node.Name.Value)
}

func (i *Interpreter) evalBinaryExpr(node *BinaryExpr) (any, error) {
	left, err := i.Evaluate(node.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.Evaluate(node.Right)
	if err != nil {
		return nil, err
	}

	return i.applyOp(node.Operator, left, right)
}

func (i *Interpreter) applyOp(op string, left, right any) (any, error) {
	// Simple type coercion logic
	// Support: int/int, float/float, string/string (+)

	switch leftVal := left.(type) {
	case int64:
		rightVal, ok := right.(int64)
		// Auto-promote right to int64 if it is int (Go implementation detail)
		if !ok {
			if rInt, ok2 := right.(int); ok2 {
				rightVal = int64(rInt)
				ok = true
			}
		}

		if ok {
			switch op {
			case "+":
				return leftVal + rightVal, nil
			case "-":
				return leftVal - rightVal, nil
			case "*":
				return leftVal * rightVal, nil
			case "/":
				if rightVal == 0 {
					return nil, fmt.Errorf("division by zero")
				}
				return leftVal / rightVal, nil
			}
		}
	case float64:
		rightVal, ok := right.(float64)
		if ok {
			switch op {
			case "+":
				return leftVal + rightVal, nil
			case "-":
				return leftVal - rightVal, nil
			case "*":
				return leftVal * rightVal, nil
			case "/":
				if rightVal == 0 {
					return nil, fmt.Errorf("division by zero")
				}
				return leftVal / rightVal, nil
			}
		}
	case string:
		rightVal, ok := right.(string)
		if ok && op == "+" {
			return leftVal + rightVal, nil
		}
	}

	return nil, fmt.Errorf("type mismatch or unknown operator: %v %s %v", left, op, right)
}

func (i *Interpreter) evalUnaryExpr(node *UnaryExpr) (any, error) {
	right, err := i.Evaluate(node.Right)
	if err != nil {
		return nil, err
	}

	switch node.Operator {
	case "-":
		switch val := right.(type) {
		case int64:
			return -val, nil
		case float64:
			return -val, nil
		}
	case "!":
		if val, ok := right.(bool); ok {
			return !val, nil
		}
	}
	return nil, fmt.Errorf("unknown unary operator: %s %v", node.Operator, right)
}

func (i *Interpreter) evalCallExpr(node *CallExpr) (any, error) {
	// Evaluate arguments
	args := []any{}
	for _, a := range node.Arguments {
		val, err := i.Evaluate(a)
		if err != nil {
			return nil, err
		}
		args = append(args, val)
	}

	// Function identifier
	ident, ok := node.Function.(*Identifier)
	if !ok {
		return nil, fmt.Errorf("function not identifier")
	}

	if fn, ok := i.Builtins[ident.Value]; ok {
		val := fn(args, i)
		return val, nil
	}

	return nil, fmt.Errorf("function not found: %s", ident.Value)
}

// Helpers

func varTypeFromToken(t TokenType) int {
	switch t {
	case TYPE_INT:
		return int(types.Int)
	case TYPE_UINT:
		return int(types.Uint)
	case TYPE_FLOAT:
		return int(types.Float)
	case TYPE_UNOFLOAT:
		return int(types.UnoFloat)
	case TYPE_BOOL:
		return int(types.Bool)
	case TYPE_STRING:
		return int(types.String)
	default:
		return int(types.Unknown)
	}
}

func (i *Interpreter) randomValue(t TokenType) any {
	switch t {
	case TYPE_INT:
		return int64(i.Rand.Intn(2000) - 1000) // Default range -1000 to 1000
	case TYPE_UINT:
		return uint64(i.Rand.Intn(2000))
	case TYPE_BOOL:
		return i.Rand.Intn(2) == 0
	case TYPE_FLOAT:
		return i.Rand.Float64() * 1000.0
	case TYPE_UNOFLOAT:
		return i.Rand.Float64()
	case TYPE_STRING:
		return i.GenerateRandomString(i.Config.RandomStringLength, i.Config.RandomStringCharset)
	}
	return nil
}

func (i *Interpreter) randomValueInRange(t TokenType, min, max any) (any, error) {
	switch t {
	case TYPE_INT:
		minVal, ok1 := toInt64(min)
		maxVal, ok2 := toInt64(max)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("invalid types for int range")
		}

		if err := checkRange(minVal, maxVal); err != nil {
			return nil, err
		}

		return i.Rand.Int63n(maxVal-minVal) + minVal, nil

	case TYPE_FLOAT:
		minVal, ok1 := toFloat64(min)
		maxVal, ok2 := toFloat64(max)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("invalid types for float range")
		}

		if err := checkRange(minVal, maxVal); err != nil {
			return nil, err
		}

		return i.Rand.Float64()*(maxVal-minVal) + minVal, nil
	}
	return nil, nil
}

func toInt64(v any) (int64, bool) {
	switch val := v.(type) {
	case int:
		return int64(val), true
	case int64:
		return val, true
	case float64:
		return int64(val), true
	}
	return 0, false
}

func toFloat64(v any) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case float64:
		return val, true
	}
	return 0, false
}

func checkRange[T int64 | float64](min, max T) error {
	if min > max {
		return fmt.Errorf("min is greater than max")
	}
	if min == max {
		return fmt.Errorf("min is equal to max")
	}
	return nil
}
