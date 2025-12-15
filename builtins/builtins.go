package builtins

import (
	"fmt"
	"wtf-script/types"
)

func RegisterBuiltins(register func(name string, fn types.IBuiltinFunc)) {
	register("print", func(args []interface{}, i types.IInterpreter) any {
		if len(args) < 1 {
			fmt.Println("print expects at least 1 argument")
			return nil
		}

		for _, arg := range args {
			fmt.Printf("%v ", arg)
		}

		fmt.Println()
		return nil
	})

	register("seed", func(args []interface{}, i types.IInterpreter) any {
		if len(args) != 1 {
			fmt.Println("seed expects exactly 1 argument")
			return nil
		}

		// Argument should be an integer
		switch v := args[0].(type) {
		case int:
			i.SetSeed(int64(v))
		case int64:
			i.SetSeed(v)
		default:
			fmt.Printf("seed expects an integer, got %T\n", args[0])
		}
		return nil
	})
}
