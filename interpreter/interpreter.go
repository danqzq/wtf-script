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

func (i *Interpreter) LogError(format string, args ...any) {
	LogError(format, args...)
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
			LogError("%s", msg)
		}
		return
	}

	_, err := i.Evaluate(program)
	if err != nil {
		LogError("%s", err)
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
	case *BlockStmt:
		return i.evalBlockStmt(node)
	case *IfStmt:
		return i.evalIfStmt(node)

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

func isLiteral(node Node) bool {
	switch node.(type) {
	case *IntegerLiteral, *FloatLiteral, *StringLiteral, *BooleanLiteral:
		return true
	}
	return false
}

func isIdentifier(node Node) bool {
	_, ok := node.(*Identifier)
	return ok
}

// validateUnofloatAssignment validates assignment to unofloat variables.
// shouldValidateStrict is true for literals and variables, false for computed expressions.
// Returns error if validation fails, or modified value for int64->float64 coercion.
func (i *Interpreter) validateUnofloatAssignment(value any, shouldValidateStrict bool, pos *Position) (any, error) {
	if val, ok := value.(float64); ok && shouldValidateStrict && (val < 0 || val > 1) {
		return nil, NewInvalidUnofloatAssignmentError(pos, val)
	}
	if val, ok := value.(types.UnofloatType); ok && shouldValidateStrict && (float64(val) < 0 || float64(val) > 1) {
		return nil, NewInvalidUnofloatAssignmentError(pos, float64(val))
	}
	if val, ok := value.(int64); ok && (val < 0 || val > 1) {
		if shouldValidateStrict {
			return nil, NewInvalidUnofloatAssignmentError(pos, float64(val))
		}
		return float64(val), nil // Coerce for clamping
	}
	return value, nil
}

// validateUintAssignment validates assignment to uint variables.
// shouldValidateStrict is true for literals and variables, false for computed expressions.
// Returns error if validation fails, or modified value for negative computed values.
func (i *Interpreter) validateUintAssignment(value any, shouldValidateStrict bool, pos *Position) (any, error) {
	if val, ok := value.(int64); ok && val < 0 {
		if shouldValidateStrict {
			return nil, NewNegativeUintAssignmentError(pos, val)
		}
		return uint64(val), nil // Allow underflow for computed values
	}
	if val, ok := value.(float64); ok && val < 0 {
		if shouldValidateStrict {
			return nil, NewNegativeUintAssignmentError(pos, int64(val))
		}
		return uint64(val), nil // Allow underflow for computed values
	}
	return value, nil
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
		pos := &Position{Line: node.Token.Line, Column: node.Token.Column}
		val, errVal = i.randomValueInRange(node.Type, minVal, maxVal, pos)
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
		shouldValidateStrict := isLiteral(node.Value) || isIdentifier(node.Value)

		// Special handling for unofloat and uint assignment validation
		pos := &Position{Line: node.Token.Line, Column: node.Token.Column}
		switch expectedType {
		case types.Unofloat:
			validatedVal, err := i.validateUnofloatAssignment(evaluated, shouldValidateStrict, pos)
			if err != nil {
				return nil, err
			}
			evaluated = validatedVal
		case types.Uint:
			validatedVal, err := i.validateUintAssignment(evaluated, shouldValidateStrict, pos)
			if err != nil {
				return nil, err
			}
			evaluated = validatedVal
		}

		err = i.checkTypeCompatibility(expectedType, evaluated, pos)
		if err != nil {
			return nil, err
		}

		val = castToType(expectedType, evaluated)
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
		pos := &Position{Line: node.Token.Line, Column: node.Token.Column}
		err = i.checkTypeCompatibility(v.Type, val, pos)
		if err != nil {
			return nil, err
		}

		shouldValidateStrict := isLiteral(node.Value) || isIdentifier(node.Value)

		// Special handling for unofloat and uint assignment validation
		switch v.Type {
		case types.Unofloat:
			validatedVal, err := i.validateUnofloatAssignment(val, shouldValidateStrict, pos)
			if err != nil {
				return nil, err
			}
			val = validatedVal
		case types.Uint:
			validatedVal, err := i.validateUintAssignment(val, shouldValidateStrict, pos)
			if err != nil {
				return nil, err
			}
			val = validatedVal
		}

		v.Value = castToType(v.Type, val)
		i.Variables[node.Name.Value] = v
		return val, nil
	}
	return nil, NewVariableNotDefinedError(node.Name)
}

func (i *Interpreter) evalBinaryExpr(node *BinaryExpr) (any, error) {
	// Handle logical operators with short-circuit evaluation
	if node.Operator == AND || node.Operator == OR {
		left, err := i.Evaluate(node.Left)
		if err != nil {
			return nil, err
		}

		leftBool := i.isTruthy(left)

		// Short-circuit: don't evaluate right if we already know the result
		if node.Operator == AND && !leftBool {
			return false, nil
		}
		if node.Operator == OR && leftBool {
			return true, nil
		}

		// Only evaluate right side if necessary
		right, err := i.Evaluate(node.Right)
		if err != nil {
			return nil, err
		}

		return i.isTruthy(right), nil
	}

	// For all other operators, evaluate both sides
	left, err := i.Evaluate(node.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.Evaluate(node.Right)
	if err != nil {
		return nil, err
	}

	pos := &Position{Line: node.Token.Line, Column: node.Token.Column}
	return i.applyOp(node.Operator, left, right, pos)
}

func defaultApplyOp[T int64 | uint64 | float64](op TokenType, l, r T) (any, error) {
	switch op {
	case PLUS:
		return l + r, nil
	case MINUS:
		return l - r, nil
	case ASTERISK:
		return l * r, nil
	case SLASH:
		if r == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return l / r, nil
	}
	return nil, fmt.Errorf("unknown operator: %s", op)
}

func (i *Interpreter) applyOp(op TokenType, left, right any, pos *Position) (any, error) {
	// Handle comparison operators separately
	if op == EQ || op == NEQ || op == LT || op == LTE || op == GT || op == GTE {
		return i.applyComparisonOp(op, left, right, pos)
	}

	leftVal, rightVal, err := i.coerceValues(left, right, pos)
	if err != nil {
		return nil, err
	}

	switch l := leftVal.(type) {
	case int64:
		r := rightVal.(int64)
		return defaultApplyOp(op, l, r)

	case uint64:
		r := rightVal.(uint64)
		return defaultApplyOp(op, l, r)

	case float64:
		r := rightVal.(float64)
		return defaultApplyOp(op, l, r)

	case types.UnofloatType:
		r := rightVal.(float64)
		switch op {
		case PLUS:
			return clampUnofloat(float64(l) + r), nil
		case MINUS:
			return clampUnofloat(float64(l) - r), nil
		case ASTERISK:
			return clampUnofloat(float64(l) * r), nil
		case SLASH:
			if r == 0.0 {
				return nil, NewDivisionByZeroError(pos)
			}
			return clampUnofloat(float64(l) / r), nil
		}

	case string:
		var r string
		var ok bool
		if r, ok = rightVal.(string); !ok {
			return nil, NewTypeMismatchError(pos, left, right)
		}
		switch op {
		case PLUS:
			return l + r, nil
		}
		return nil, NewRuntimeError(pos, "unknown string operator: %s", op)
	}

	return nil, NewUnknownOperatorError(pos, op, left, right)
}

func defaultComparisonOp[T int64 | uint64 | float64 | string](op TokenType, l, r T) (bool, error) {
	switch op {
	case EQ:
		return l == r, nil
	case NEQ:
		return l != r, nil
	case LT:
		return l < r, nil
	case LTE:
		return l <= r, nil
	case GT:
		return l > r, nil
	case GTE:
		return l >= r, nil
	}
	return false, fmt.Errorf("unknown comparison operator: %s", op)
}

func (i *Interpreter) applyComparisonOp(op TokenType, left, right any, pos *Position) (any, error) {
	leftVal, rightVal, err := i.coerceValues(left, right, pos)
	if err != nil {
		return nil, err
	}

	switch l := leftVal.(type) {
	case int64:
		r := rightVal.(int64)
		return defaultComparisonOp(op, l, r)
	case uint64:
		r := rightVal.(uint64)
		return defaultComparisonOp(op, l, r)
	case float64:
		r := rightVal.(float64)
		return defaultComparisonOp(op, l, r)
	case types.UnofloatType:
		r := rightVal.(float64)
		fl := float64(l)
		return defaultComparisonOp(op, fl, r)
	case string:
		r := rightVal.(string)
		return defaultComparisonOp(op, l, r)
	case bool:
		r, ok := rightVal.(bool)
		if !ok {
			return nil, NewRuntimeError(pos, "cannot compare bool with %T", rightVal)
		}
		switch op {
		case EQ:
			return l == r, nil
		case NEQ:
			return l != r, nil
		default:
			return nil, NewRuntimeError(pos, "operator %s not supported for bool", op)
		}
	}

	return nil, NewUnknownOperatorError(pos, op, left, right)
}
func (i *Interpreter) isTruthy(val any) bool {
	switch v := val.(type) {
	case bool:
		return v
	case int64:
		return v != 0
	case uint64:
		return v != 0
	case float64:
		return v != 0.0
	case types.UnofloatType:
		return float64(v) != 0.0
	case string:
		return len(v) > 0
	case nil:
		return false
	default:
		return true
	}
}

func coerceHandleRight[T int64 | uint64 | float64 | types.UnofloatType](l T, r any, pos *Position) (T, T, error) {
	switch rv := r.(type) {
	case int64:
		return l, T(rv), nil
	case uint64:
		return l, T(rv), nil
	case float64:
		return l, T(rv), nil
	case types.UnofloatType:
		return l, T(float64(rv)), nil
	default:
		var zero T
		return zero, zero, NewTypeMismatchError(pos, l, r)
	}
}

func (i *Interpreter) coerceValues(left, right any, pos *Position) (any, any, error) {
	switch l := left.(type) {
	case int64:
		// Left is Int: Coerce Right to Int (FCFS)
		return coerceHandleRight(l, right, pos)
	case uint64:
		// Left is Uint: Coerce Right to Uint (FCFS)
		return coerceHandleRight(l, right, pos)
	case float64:
		// Left is Float: Coerce Right to Float (FCFS)
		return coerceHandleRight(l, right, pos)
	case types.UnofloatType:
		// Left is Unofloat: Coerce Right to Float, keep left as Unofloat (FCFS)
		switch r := right.(type) {
		case types.UnofloatType:
			return l, float64(r), nil
		case float64:
			return l, r, nil
		case int64:
			return l, float64(r), nil
		case uint64:
			return l, float64(r), nil
		default:
			return nil, nil, i.typeMismatchError(left, right, pos)
		}

	case string:
		if r, ok := right.(string); ok {
			return l, r, nil
		}
		return nil, nil, i.typeMismatchError(left, right, pos)

	case bool:
		if r, ok := right.(bool); ok {
			return l, r, nil
		}
		return nil, nil, i.typeMismatchError(left, right, pos)
	}
	return nil, nil, NewRuntimeError(pos, "unsupported left operand type: %T", left)
}

func (i *Interpreter) typeMismatchError(left, right any, pos *Position) error {
	return NewTypeMismatchError(pos, left, right)
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
		case uint64:
			return 0 - val, nil
		case float64:
			return -val, nil
		case types.UnofloatType:
			return -float64(val), nil
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

func (i *Interpreter) evalBlockStmt(block *BlockStmt) (any, error) {
	var result any
	for _, statement := range block.Statements {
		val, err := i.Evaluate(statement)
		if err != nil {
			return nil, err
		}
		result = val
	}
	return result, nil
}

func (i *Interpreter) evalIfStmt(node *IfStmt) (any, error) {
	var condition bool
	pos := &Position{Line: node.Token.Line, Column: node.Token.Column}
	if node.Token.Type == IFRAND {
		// ifrand statement
		if node.Condition != nil {
			// ifrand(probability)
			probVal, err := i.Evaluate(node.Condition)
			if err != nil {
				return nil, err
			}

			var probability float64
			switch p := probVal.(type) {
			case float64:
				probability = p
			case int64:
				probability = float64(p)
			case uint64:
				probability = float64(p)
			case types.UnofloatType:
				probability = float64(p)
			default:
				return nil, NewRuntimeError(pos, "ifrand probability must be a number, got %T", probVal)
			}

			if probability < UnofloatMin || probability > UnofloatMax {
				return nil, NewRuntimeError(pos, "ifrand probability must be between 0 and 1, got %f", probability)
			}

			condition = i.Rand.Float64() < probability
		} else {
			// ifrand without probability (default 0.5)
			condition = i.Rand.Float64() < DefaultIfrandProbability
		}
	} else {
		// Regular if statement
		condVal, err := i.Evaluate(node.Condition)
		if err != nil {
			return nil, err
		}

		boolCond, ok := condVal.(bool)
		if !ok {
			return nil, NewRuntimeError(pos, "if condition must evaluate to bool, got %T", condVal)
		}
		condition = boolCond
	}

	if condition {
		return i.Evaluate(node.Consequence)
	} else if node.Alternative != nil {
		return i.Evaluate(node.Alternative)
	}

	return nil, nil
}
