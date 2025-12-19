package main

import (
	"flag"
	"os"
	"wtf-script/config"
	"wtf-script/interpreter"
)

func main() {
	configFile := flag.String("config", "", "Path to JSON configuration file")
	flag.Parse()

	if flag.NArg() < 1 {
		interpreter.LogError("Usage: wtf [--config <config.json>] <file.wtf>")
		return
	}

	filename := flag.Arg(0)
	content, err := os.ReadFile(filename)
	if err != nil {
		interpreter.LogError("Error reading file: %v", err)
		return
	}

	// Load config if provided
	var cfg *config.Config
	if *configFile != "" {
		cfg, err = config.LoadConfigFromFile(*configFile)
		if err != nil {
			interpreter.LogError("Error loading config: %v", err)
			return
		}
	}

	i := interpreter.NewInterpreter(cfg)
	i.Execute(string(content))
}
