package main

import (
	"fmt"
	"os"
	"wtf-script/interpreter"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: wtf <file.wtf>")
		os.Exit(1)
	}

	filename := os.Args[1]
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	i := interpreter.NewInterpreter()
	i.Run(string(content))
}
