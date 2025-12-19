package interpreter

import "fmt"

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

// LogError prints an error message in red
func LogError(format string, args ...any) {
	fmt.Printf(ColorRed+format+ColorReset+"\n", args...)
}

// LogInfo prints an informational message
func LogInfo(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}

// LogWarning prints a warning message in yellow
func LogWarning(format string, args ...any) {
	fmt.Printf(ColorYellow+format+ColorReset+"\n", args...)
}

// LogSuccess prints a success message in green
func LogSuccess(format string, args ...any) {
	fmt.Printf(ColorGreen+format+ColorReset+"\n", args...)
}
