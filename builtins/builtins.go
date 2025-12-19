package builtins

import (
	"fmt"
	"wtf-script/types"
)

const (
	PRINT  = "print"
	SEED   = "seed"
	TYPEOF = "typeof"
)

func RegisterBuiltins(register func(name string, fn types.IBuiltinFunc)) {
	register(PRINT, func(args []any, i types.IInterpreter) any {
		if len(args) < 1 {
			i.LogError("print expects at least 1 argument")
			return nil
		}

		for _, arg := range args {
			switch v := arg.(type) {
			case float64:
				fmt.Printf("%f ", v)
			case types.UnofloatType:
				fmt.Printf("%f ", float64(v))
			default:
				fmt.Printf("%v ", arg)
			}
		}

		fmt.Println()
		return nil
	})

	register(SEED, func(args []any, i types.IInterpreter) any {
		if len(args) != 1 {
			i.LogError("seed expects exactly 1 argument")
			return nil
		}

		// Argument should be an integer
		switch v := args[0].(type) {
		case int:
			i.SetSeed(int64(v))
		case int64:
			i.SetSeed(v)
		default:
			i.LogError("seed expects an integer, got %T", args[0])
		}
		return nil
	})

	register(TYPEOF, func(args []any, i types.IInterpreter) any {
		if len(args) != 1 {
			i.LogError("typeof expects exactly 1 argument")
			return nil
		}

		switch args[0].(type) {
		case int64:
			return "int"
		case uint64:
			return "uint"
		case float64:
			return "float"
		case types.UnofloatType:
			return "unofloat"
		case string:
			return "string"
		case bool:
			return "bool"
		case nil:
			return "nil"
		default:
			return "unknown"
		}
	})
}
