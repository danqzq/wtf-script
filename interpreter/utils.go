package interpreter

import (
	"strings"
	"unicode"
)

func toFloat64(val interface{}) float64 {
	switch v := val.(type) {
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case float64:
		return v
	default:
		return 0
	}
}

func isQuotedString(s string) bool {
	return strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")
}

func splitArgsPreserveStrings(s string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false

	for i := 0; i < len(s); i++ {
		c := s[i]

		if c == '"' {
			inQuotes = !inQuotes
			current.WriteByte(c)
		} else if c == ',' && !inQuotes {
			args = append(args, strings.TrimSpace(current.String()))
			current.Reset()
		} else {
			current.WriteByte(c)
		}
	}

	if current.Len() > 0 {
		args = append(args, strings.TrimSpace(current.String()))
	}

	return args
}

func isFunctionCall(line string) bool {
	line = strings.TrimSpace(line)
	if len(line) < 3 {
		return false
	}
	openParen := strings.Index(line, "(")
	closeParen := strings.LastIndex(line, ");")
	if openParen <= 0 || closeParen != len(line)-2 {
		return false
	}
	funcName := line[:openParen]
	return isValidIdentifier(funcName)
}

func isValidIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}
	if !unicode.IsLetter(rune(s[0])) && s[0] != '_' {
		return false
	}
	for _, c := range s {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c != '_' {
			return false
		}
	}
	return true
}
