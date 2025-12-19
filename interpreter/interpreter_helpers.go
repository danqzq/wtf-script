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
		return int(types.Unofloat)
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
		rangeSize := i.Config.Int.Max - i.Config.Int.Min
		return i.Config.Int.Min + i.Rand.Int63n(rangeSize) + RangeInclusiveOffset
	case TYPE_UINT:
		rangeSize := i.Config.Uint.Max - i.Config.Uint.Min
		return i.Config.Uint.Min + uint64(i.Rand.Int63n(int64(rangeSize))) + RangeInclusiveOffset
	case TYPE_BOOL:
		return i.Rand.Intn(RandomBoolChoices) == 0
	case TYPE_FLOAT:
		rangeSize := i.Config.Float.Max - i.Config.Float.Min
		return i.Config.Float.Min + i.Rand.Float64()*rangeSize
	case TYPE_UNOFLOAT:
		rangeSize := i.Config.Unofloat.Max - i.Config.Unofloat.Min
		return types.UnofloatType(i.Config.Unofloat.Min + i.Rand.Float64()*rangeSize)
	case TYPE_STRING:
		return i.GenerateRandomString(int(i.Config.Length.Min), i.Config.Charset)
	}
	return nil
}

func (i *Interpreter) randomValueInRange(t TokenType, min, max any, pos *Position) (any, error) {
	switch t {
	case TYPE_INT:
		minVal, ok1 := toInt64(min)
		maxVal, ok2 := toInt64(max)
		if !ok1 || !ok2 {
			return nil, NewRuntimeError(pos, "invalid types for int range")
		}

		if err := checkRange(minVal, maxVal, pos); err != nil {
			return nil, err
		}

		return i.Rand.Int63n(maxVal-minVal) + minVal, nil

	case TYPE_FLOAT:
		minVal, ok1 := toFloat64(min)
		maxVal, ok2 := toFloat64(max)
		if !ok1 || !ok2 {
			return nil, NewRuntimeError(pos, "invalid types for float range")
		}

		if err := checkRange(minVal, maxVal, pos); err != nil {
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
	case types.UnofloatType:
		return float64(val), true
	}
	return 0, false
}

func checkRange[T int64 | float64](min, max T, pos *Position) error {
	if min > max {
		return NewInvalidRangeError(pos, "min is greater than max")
	}
	if min == max {
		return NewInvalidRangeError(pos, "min is equal to max")
	}
	return nil
}

func getTypeString(value any) string {
	switch value.(type) {
	case int64:
		return "int"
	case uint64:
		return "uint"
	case float64:
		return "float"
	case types.UnofloatType:
		return "unofloat"
	case bool:
		return "bool"
	case string:
		return "string"
	default:
		return "unknown"
	}
}

func defaultTypeCompatibility(expectedType *types.VarType, value any, pos *Position) error {
	switch value.(type) {
	case int64, uint64, float64, types.UnofloatType:
		return nil
	}
	return NewRuntimeError(pos, "type mistmatch: expected %s, got %s", expectedType.String(), getTypeString(value))
}

func (i *Interpreter) checkTypeCompatibility(expectedType types.VarType, value any, pos *Position) error {
	switch expectedType {
	case types.Int:
		return defaultTypeCompatibility(&expectedType, value, pos)
	case types.Uint:
		return defaultTypeCompatibility(&expectedType, value, pos)
	case types.Float:
		return defaultTypeCompatibility(&expectedType, value, pos)
	case types.Unofloat:
		return defaultTypeCompatibility(&expectedType, value, pos)
	case types.Bool:
		if _, ok := value.(bool); !ok {
			return NewRuntimeError(pos, "type mismatch: expected bool, got %T", value)
		}
	case types.String:
		if _, ok := value.(string); !ok {
			return NewRuntimeError(pos, "type mismatch: expected string, got %T", value)
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
	case types.Unofloat:
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

func clampUnofloat(value float64) types.UnofloatType {
	return types.UnofloatType(math.Max(UnofloatMin, math.Min(UnofloatMax, value)))
}
