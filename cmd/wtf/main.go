package main

import (
	"flag"
	"fmt"
	"os"
	"wtf-script/config"
	"wtf-script/interpreter"
)

func main() {
	configFile := flag.String("config", "", "Path to JSON configuration file")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: wtf [--config <config.json>] <file.wtf>")
		os.Exit(1)
	}

	filename := flag.Arg(0)
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	// Load config if provided
	var cfg *config.Config
	if *configFile != "" {
		cfg, err = config.LoadConfigFromFile(*configFile)
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}
	}

	i := interpreter.NewInterpreter(cfg)
	i.Execute(string(content))
}
