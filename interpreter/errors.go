package interpreter

import "fmt"

func PrintError(line, column int, errorType, msg string) string {
	return fmt.Sprintf("[Line %d, Col %d] %s error: %s", line, column, errorType, msg)
}
