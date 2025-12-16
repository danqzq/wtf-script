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

func (i *Interpreter) GetConfig() *config.Config {
	return i.Config
}

func (i *Interpreter) GenerateRandomString(n int, charset string) string {
	b := make([]byte, n)
	for j := range b {
		b[j] = charset[i.Rand.Intn(len(charset))]
	}
	return string(b)
}

func (i *Interpreter) SetSeed(seed int64) {
	i.Rand.Seed(seed)
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
	return nil, NewIdentifierNotFoundError(node)
}

func (i *Interpreter) evalVarDecl(node *VarDecl) (any, error) {
	var val any

	if node.RangeMin != nil && node.RangeMax != nil {
		// Handles: int(0, 100) x;
		minVal, err := i.Evaluate(node.RangeMin)
		if err != nil {
			return nil, err
		}
		maxVal, err := i.Evaluate(node.RangeMax)
		if err != nil {
			return nil, err
		}

		var errVal error
		val, errVal = i.randomValueInRange(node.Type, minVal, maxVal, node.Token.Line, node.Token.Column)
		if errVal != nil {
			return nil, errVal
		}
	} else if node.Value != nil {
		// Handles: int x = 5;
		evaluated, err := i.Evaluate(node.Value)
		if err != nil {
			return nil, err
		}

		expectedType := types.VarType(varTypeFromToken(node.Type))
		err = i.checkTypeCompatibility(expectedType, evaluated, node.Token.Line, node.Token.Column)
		if err != nil {
			return nil, err
		}

		val = evaluated
	} else {
		// Handles: int x; (random default)
		val = i.randomValue(node.Type)
	}

	i.Variables[node.Name.Value] = types.Variable{
		Type:  types.VarType(varTypeFromToken(node.Type)),
		Value: val,
	}
	return val, nil
}

func (i *Interpreter) evalAssignStmt(node *AssignStmt) (any, error) {
	val, err := i.Evaluate(node.Value)
	if err != nil {
		return nil, err
	}

	if v, ok := i.Variables[node.Name.Value]; ok {
		err = i.checkTypeCompatibility(v.Type, val, node.Token.Line, node.Token.Column)
		if err != nil {
			return nil, err
		}
		v.Value = val
		i.Variables[node.Name.Value] = v
		return val, nil
	}
	return nil, NewVariableNotDefinedError(node.Name)
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

	return i.applyOp(node.Operator, left, right, node.Token.Line, node.Token.Column)
}

func (i *Interpreter) applyOp(op string, left, right any, line, col int) (any, error) {
	leftVal, rightVal, err := i.coerceValues(left, right, line, col)
	if err != nil {
		return nil, err
	}

	switch l := leftVal.(type) {
	case int64:
		r := rightVal.(int64)
		switch op {
		case "+":
			return l + r, nil
		case "-":
			return l - r, nil
		case "*":
			return l * r, nil
		case "/":
			if r == 0 {
				return nil, NewDivisionByZeroError(line, col)
			}
			return l / r, nil
		}

	case float64:
		r := rightVal.(float64)
		switch op {
		case "+":
			return l + r, nil
		case "-":
			return l - r, nil
		case "*":
			return l * r, nil
		case "/":
			if r == 0.0 {
				return nil, NewDivisionByZeroError(line, col)
			}
			return l / r, nil
		}

	case string:
		if op != "+" {
			return nil, NewRuntimeError(line, col, "unknown string operator: %s", op)
		}
		return l + rightVal.(string), nil
	}

	return nil, NewUnknownOperatorError(line, col, op, left, right)
}

func (i *Interpreter) coerceValues(left, right any, line, col int) (any, any, error) {
	switch l := left.(type) {
	case int64:
		// Left is Int: Coerce Right to Int (FCFS)
		switch r := right.(type) {
		case int64:
			return l, r, nil
		case float64:
			return l, int64(r), nil // Coerce float to int
		default:
			return nil, nil, i.typeMismatchError(left, right, line, col)
		}

	case float64:
		// Left is Float: Coerce Right to Float (FCFS)
		switch r := right.(type) {
		case float64:
			return l, r, nil
		case int64:
			return l, float64(r), nil // Coerce int to float
		default:
			return nil, nil, i.typeMismatchError(left, right, line, col)
		}

	case string:
		if r, ok := right.(string); ok {
			return l, r, nil
		}
		return nil, nil, i.typeMismatchError(left, right, line, col)

	case bool:
		// Bool does not support arithmetic/concatenation with other types
		return nil, nil, i.typeMismatchError(left, right, line, col)
	}
	return nil, nil, NewRuntimeError(line, col, "unsupported left operand type: %T", left)
}

func (i *Interpreter) typeMismatchError(left, right any, line, col int) error {
	return NewTypeMismatchError(line, col, left, right)
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
	return nil, NewUnknownUnaryOperatorError(node, right)
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
		return nil, NewInvalidFunctionCallError(node, "function expression must be an identifier")
	}

	if fn, ok := i.Builtins[ident.Value]; ok {
		val := fn(args, i)
		return val, nil
	}

	return nil, NewFunctionNotFoundError(node, ident.Value)
}
