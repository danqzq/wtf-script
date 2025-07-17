package interpreter

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"wtf-script/types"
)

func parseVarType(typeToken string) types.VarType {
	if strings.HasPrefix(typeToken, "int(") {
		return types.Int
	}
	if strings.HasPrefix(typeToken, "float(") {
		return types.Float
	}
	switch typeToken {
	case "int":
		return types.Int
	case "uint":
		return types.Uint
	case "float":
		return types.Float
	case "unofloat":
		return types.UnoFloat
	case "bool":
		return types.Bool
	case "string":
		return types.String
	default:
		return types.Unknown
	}
}

func (i *Interpreter) parseLine(line string) {
	line = strings.TrimSpace(line)

	// IGNORE COMMENTS AND EMPTY LINES
	if line == "" || strings.HasPrefix(line, "//") {
		return
	}
	if strings.Contains(line, "//") {
		line = strings.Split(line, "//")[0]
		line = strings.TrimSpace(line)
		if line == "" {
			return
		}
	}

	// BUILT-IN FUNCTION CALLS
	if isFunctionCall(line) {
		i.handleFunctionCall(line)
		return
	}

	// ASSIGNMENT
	if strings.Contains(line, "=") {
		parts := strings.SplitN(line, "=", 2)
		left := strings.TrimSpace(parts[0])
		right := strings.TrimSpace(strings.TrimSuffix(parts[1], ";"))

		tokens := strings.Fields(left)
		if len(tokens) == 2 {
			typeToken := tokens[0]
			varName := tokens[1]
			varType := parseVarType(typeToken)
			err := i.validateExpressionOrLiteral(varType, right)
			if err != nil {
				fmt.Println(err)
				return
			}

			i.createOrReplaceVariable(varType, right, varName)
			return
		}

		// ASSIGNMENT TO EXISTING VARIABLE
		varName := left
		if v, ok := i.Variables[varName]; ok {
			err := i.validateExpressionOrLiteral(v.Type, right)
			if err != nil {
				fmt.Println(err)
				return
			}

			i.createOrReplaceVariable(v.Type, right, varName)
			return
		}

		fmt.Println("Undefined variable:", varName)
		return
	}

	if strings.HasSuffix(line, ";") {
		decl := strings.TrimSuffix(line, ";")
		tokens := strings.Fields(decl)

		if len(tokens) >= 2 {
			typeToken := tokens[0]
			varName := tokens[1]

			if strings.Contains(typeToken, "(") && !strings.Contains(typeToken, ")") {
				idx := 1
				for idx < len(tokens) && !strings.Contains(typeToken, ")") {
					typeToken += tokens[idx]
					idx++
				}

				if idx < len(tokens) {
					varName = tokens[idx]
				} else {
					fmt.Println("Invalid declaration syntax")
					return
				}
			}

			var varType types.VarType
			var val interface{}
			defer func() {
				i.Variables[varName] = types.Variable{Type: varType, Value: val}
			}()

			if strings.Contains(typeToken, "(") && strings.Contains(typeToken, ")") {
				baseType := typeToken[:strings.Index(typeToken, "(")]
				rangePart := typeToken[strings.Index(typeToken, "(")+1 : strings.Index(typeToken, ")")]
				ranges := strings.Split(rangePart, ",")
				if len(ranges) != 2 {
					fmt.Println("Invalid range specification")
					return
				}
				minStr := strings.TrimSpace(ranges[0])
				maxStr := strings.TrimSpace(ranges[1])

				varType = parseVarType(baseType)
				minVal := i.parseLiteral(varType, minStr)
				maxVal := i.parseLiteral(varType, maxStr)
				val = i.randomValueInRange(varType, minVal, maxVal)
				return
			}

			varType = parseVarType(typeToken)
			val = i.randomValue(varType)
			return
		}
	}

	fmt.Println("Unrecognized line:", line)
}

func (i *Interpreter) validateExpressionOrLiteral(varType types.VarType, expr string) error {
	if strings.Contains(expr, "(") || strings.Contains(expr, ")") {
		return nil
	}
	if strings.Contains(expr, "+") || strings.Contains(expr[1:], "-") ||
		strings.Contains(expr, "*") || strings.Contains(expr, "/") {
		return nil
	}
	switch varType {
	case types.Int:
		_, err := strconv.ParseInt(expr, 10, 64)
		if err != nil {
			return errors.New("Failed to parse int literal: " + expr)
		}
		return nil
	case types.Uint:
		_, err := strconv.ParseUint(expr, 10, 64)
		if err != nil {
			return errors.New("Failed to parse uint literal: " + expr)
		}
		return nil
	case types.Float, types.UnoFloat:
		_, err := strconv.ParseFloat(expr, 64)
		if err != nil {
			return errors.New("Failed to parse float/unofloat literal: " + expr)
		}
		return nil
	case types.Bool:
		_, err := strconv.ParseBool(expr)
		if err != nil {
			return errors.New("Failed to parse bool literal: " + expr)
		}
		return nil
	case types.String:
		if isQuotedString(expr) {
			return nil
		}
		return errors.New("String literal must be quoted: " + expr)
	default:
		return nil
	}
}

func (i *Interpreter) createOrReplaceVariable(varType types.VarType, right string, varName string) {
	var val interface{}

	defer func() {
		i.Variables[varName] = types.Variable{Type: varType, Value: val}
	}()

	if varType == types.String && isQuotedString(right) {
		val = right[1 : len(right)-1]
		return
	}

	if varType == types.Bool {
		if right == trueValue {
			val = true
			return
		} else if right == falseValue {
			val = false
			return
		}
		fmt.Println("Invalid boolean value:", right)
		return
	}

	t := tokenize(right)
	pos := 0
	val = i.parseExpression(t, &pos)

	switch varType {
	case types.Int:
		val = int(val.(float64))
	case types.Uint:
		val = uint64(val.(float64))
	case types.Float:
		val = val.(float64)
	case types.UnoFloat:
		val = val.(float64)
	}
}

func (i *Interpreter) handleTypeDeclaration(line string) {
	if strings.Contains(line, "(") && strings.Contains(line, ")") && !strings.Contains(line, "=") {
		typePart := line[:strings.Index(line, ")")+1]
		rest := strings.TrimSpace(line[strings.Index(line, ")")+1:])
		varName := strings.TrimSuffix(rest, ";")

		typeName := typePart[:strings.Index(typePart, "(")]
		rangePart := typePart[strings.Index(typePart, "(")+1 : strings.Index(typePart, ")")]
		ranges := strings.Split(rangePart, ",")
		if len(ranges) != 2 {
			fmt.Println("Invalid range specification")
			return
		}
		minStr := strings.TrimSpace(ranges[0])
		maxStr := strings.TrimSpace(ranges[1])

		varType := parseVarType(typeName)

		minVal := i.parseLiteral(varType, minStr)
		maxVal := i.parseLiteral(varType, maxStr)

		val := i.randomValueInRange(varType, minVal, maxVal)
		i.Variables[varName] = types.Variable{Type: varType, Value: val}
		return
	} else {
		tokens := strings.SplitN(line, " ", 2)
		if len(tokens) < 2 {
			return
		}
		varType := parseVarType(tokens[0])
		rest := strings.TrimSpace(tokens[1])

		if strings.Contains(rest, "=") {
			assignParts := strings.SplitN(rest, "=", 2)
			varName := strings.TrimSpace(assignParts[0])
			valStr := strings.TrimSuffix(strings.TrimSpace(assignParts[1]), ";")
			val := i.parseLiteral(varType, valStr)
			i.Variables[varName] = types.Variable{Type: varType, Value: val}
		} else {
			varName := strings.TrimSuffix(rest, ";")
			val := i.randomValue(varType)
			i.Variables[varName] = types.Variable{Type: varType, Value: val}
		}
	}
}

func (i *Interpreter) handleFunctionCall(line string) {
	indexOfOpen := strings.Index(line, "(")
	indexOfClose := strings.Index(line, ")")

	fnName := line[:indexOfOpen]
	fnName = strings.TrimSpace(fnName)
	argsStr := line[indexOfOpen+1 : indexOfClose]
	args := splitArgsPreserveStrings(argsStr)

	if fn, ok := i.Builtins[fnName]; ok {
		fn(args, i)
	} else {
		fmt.Println("Unknown function:", fnName)
	}
}

func (i *Interpreter) parseLiteral(varType types.VarType, valStr string) interface{} {
	var err error
	switch varType {
	case types.Int:
		var parsed int
		_, err = fmt.Sscanf(valStr, "%d", &parsed)
		if err != nil {
			fmt.Println("Invalid integer value:", valStr)
			return nil
		}
		return parsed
	case types.Uint:
		var parsed int64
		_, err = fmt.Sscanf(valStr, "%d", &parsed)
		if err != nil {
			fmt.Println("Invalid unsigned integer value:", valStr)
			return nil
		}
		if parsed < 0 {
			fmt.Println("Uint cannot be negative:", valStr)
			return nil
		}
		return uint64(parsed)
	case types.Float:
		var parsed float64
		_, err = fmt.Sscanf(valStr, "%f", &parsed)
		if err != nil {
			fmt.Println("Invalid float value:", valStr)
			return nil
		}
		return parsed
	case types.UnoFloat:
		var parsed float64
		_, err = fmt.Sscanf(valStr, "%f", &parsed)
		if err != nil {
			fmt.Println("Invalid UnoFloat value:", valStr)
			return nil
		}
		if parsed < 0 || parsed > 1 {
			fmt.Println("UnoFloat must be between 0 and 1:", valStr)
			return nil
		}
		return parsed
	case types.Bool:
		isTrue := valStr == "true"
		if isTrue {
			return true
		}
		if valStr != "false" {
			fmt.Println("Invalid boolean value:", valStr)
			return nil
		}
		return false
	case types.String:
		return strings.Trim(valStr, "\"")
	default:
		return nil
	}
}

func (i *Interpreter) randomValue(varType types.VarType) interface{} {
	switch varType {
	case types.Int:
		return i.Rand.Intn(maxIntRange*2) - maxIntRange
	case types.Uint:
		return uint64(i.Rand.Intn(maxIntRange + 1))
	case types.Float:
		return i.Rand.Float64()*maxFloatRange*2 - maxFloatRange
	case types.UnoFloat:
		return i.Rand.Float64()
	case types.Bool:
		return i.Rand.Intn(2) == 0
	case types.String:
		length := i.Rand.Intn(maxStringLength-minStringLength+1) + minStringLength
		letters := []rune(randomStringCharacters)
		var sb strings.Builder
		for j := 0; j < length; j++ {
			sb.WriteRune(letters[i.Rand.Intn(len(letters))])
		}
		return sb.String()
	default:
		return nil
	}
}

func (i *Interpreter) randomValueInRange(varType types.VarType, min interface{}, max interface{}) interface{} {
	switch varType {
	case types.Int:
		minInt := min.(int)
		maxInt := max.(int)
		return i.Rand.Intn(maxInt-minInt+1) + minInt
	case types.Float:
		minF := min.(float64)
		maxF := max.(float64)
		return i.Rand.Float64()*(maxF-minF) + minF
	default:
		fmt.Println("Range not supported for this type")
		return nil
	}
}

func (i *Interpreter) parseExpression(tokens []string, pos *int) float64 {
	result := i.parseTerm(tokens, pos)

	for *pos < len(tokens) {
		op := tokens[*pos]
		if op != "+" && op != "-" {
			break
		}
		*pos++
		right := i.parseTerm(tokens, pos)
		if op == "+" {
			result += right
		} else {
			result -= right
		}
	}

	return result
}

func (i *Interpreter) parseTerm(tokens []string, pos *int) float64 {
	result := i.parseFactor(tokens, pos)

	for *pos < len(tokens) {
		op := tokens[*pos]
		if op != "*" && op != "/" {
			break
		}
		*pos++
		right := i.parseFactor(tokens, pos)
		if op == "*" {
			result *= right
		} else {
			if right == 0 {
				fmt.Println("Division by zero")
				return 0
			}
			result /= right
		}
	}

	return result
}

func (i *Interpreter) parseFactor(tokens []string, pos *int) float64 {
	token := tokens[*pos]
	*pos++

	if token == "-" {
		if *pos < len(tokens) {
			right := i.parseFactor(tokens, pos)
			return -right
		}
		fmt.Println("Invalid negation")
		return 0
	}

	if token == "(" {
		result := i.parseExpression(tokens, pos)
		if *pos < len(tokens) && tokens[*pos] == ")" {
			*pos++
		} else {
			fmt.Println("Missing closing parenthesis")
		}
		return result
	}

	if v, ok := i.Variables[token]; ok {
		return toFloat64(v.Value)
	}

	var num float64
	_, err := fmt.Sscanf(token, "%f", &num)
	if err != nil {
		fmt.Println("Invalid number:", token)
		return 0
	}
	return num
}
