package interpreter

import (
	"math"
	"wtf-script/types"
)

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

func (i *Interpreter) randomValueInRange(t TokenType, min, max any, line, col int) (any, error) {
	switch t {
	case TYPE_INT:
		minVal, ok1 := toInt64(min)
		maxVal, ok2 := toInt64(max)
		if !ok1 || !ok2 {
			return nil, NewRuntimeError(line, col, "invalid types for int range")
		}

		if err := checkRange(minVal, maxVal, line, col); err != nil {
			return nil, err
		}

		return i.Rand.Int63n(maxVal-minVal) + minVal, nil

	case TYPE_FLOAT:
		minVal, ok1 := toFloat64(min)
		maxVal, ok2 := toFloat64(max)
		if !ok1 || !ok2 {
			return nil, NewRuntimeError(line, col, "invalid types for float range")
		}

		if err := checkRange(minVal, maxVal, line, col); err != nil {
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

func checkRange[T int64 | float64](min, max T, line, col int) error {
	if min > max {
		return NewInvalidRangeError(line, col, "min is greater than max")
	}
	if min == max {
		return NewInvalidRangeError(line, col, "min is equal to max")
	}
	return nil
}

func (i *Interpreter) checkTypeCompatibility(expectedType types.VarType, value any, line, col int) error {
	switch expectedType {
	case types.Int:
		if _, ok := value.(int64); ok {
			return nil
		}
		if _, ok := value.(uint64); ok {
			return nil
		}
		if _, ok := value.(float64); ok {
			return nil
		}
		return NewRuntimeError(line, col, "type mismatch: expected int, got %T", value)
	case types.Uint:
		if _, ok := value.(uint64); ok {
			return nil
		}
		if _, ok := value.(int64); ok {
			return nil
		}
		return NewRuntimeError(line, col, "type mismatch: expected uint, got %T", value)
	case types.Float, types.UnoFloat:
		if _, ok := value.(float64); ok {
			return nil
		}
		if _, ok := value.(int64); ok {
			return nil
		}
		if _, ok := value.(uint64); ok {
			return nil
		}
		expectedTypeStr := "float"
		if expectedType == types.UnoFloat {
			expectedTypeStr = "unofloat"
		}
		return NewRuntimeError(line, col, "type mismatch: expected %s, got %T", expectedTypeStr, value)
	case types.Bool:
		if _, ok := value.(bool); !ok {
			return NewRuntimeError(line, col, "type mismatch: expected bool, got %T", value)
		}
	case types.String:
		if _, ok := value.(string); !ok {
			return NewRuntimeError(line, col, "type mismatch: expected string, got %T", value)
		}
	}
	return nil
}

func castToType(expectedType types.VarType, value any) any {
	switch expectedType {
	case types.Uint:
		if value, ok := value.(int64); ok {
			return uint64(value)
		}
		if value, ok := value.(float64); ok {
			return uint64(value)
		}
	case types.Float:
		if value, ok := value.(int64); ok {
			return float64(value)
		}
		if value, ok := value.(uint64); ok {
			return float64(value)
		}
	case types.UnoFloat:
		if value, ok := value.(int64); ok {
			return clampUnofloat(float64(value))
		}
		if value, ok := value.(uint64); ok {
			return clampUnofloat(float64(value))
		}
		if value, ok := value.(float64); ok {
			return clampUnofloat(value)
		}
	case types.Int:
		if value, ok := value.(uint64); ok {
			return int64(value)
		}
		if value, ok := value.(float64); ok {
			return int64(value)
		}
	}
	return value
}

func clampUnofloat(value float64) float64 {
	return math.Max(0.0, math.Min(1.0, value))
}
