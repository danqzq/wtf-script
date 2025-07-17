package builtins

import (
	"fmt"
	"strconv"
	"wtf-script/types"
)

func RegisterBuiltins(register func(name string, fn types.IBuiltinFunc)) {
	register("print", func(args []string, i types.IInterpreter) {
		if len(args) < 1 {
			fmt.Println("print expects at least 1 argument")
			return
		}

		for _, arg := range args {
			valStr := resolveArgValue(arg, i)
			fmt.Printf("%v ", valStr)
		}

		fmt.Println()
	})

	register("seed", func(args []string, i types.IInterpreter) {
		if len(args) != 1 {
			fmt.Println("seed expects exactly 1 argument")
			return
		}

		arg := args[0]
		if v, ok := i.GetVariable(arg); ok {
			if num, err := strconv.ParseInt(v.Value.(string), 10, 64); err == nil {
				i.SetSeed(num)
				return
			}
			fmt.Println("`" + arg + "` is not a valid integer variable.")
			return
		}

		if num, err := strconv.ParseInt(arg, 10, 64); err == nil {
			i.SetSeed(num)
			return
		}

		fmt.Println("`" + arg + "` is not a valid integer literal or variable.")
	})
}

func resolveArgValue(arg string, i types.IInterpreter) string {
	if len(arg) >= 2 && arg[0] == '"' && arg[len(arg)-1] == '"' {
		return arg[1 : len(arg)-1]
	}

	if v, ok := i.GetVariable(arg); ok {
		return v.ToString()
	}

	if _, err := strconv.ParseInt(arg, 10, 64); err == nil {
		return arg
	}

	if _, err := strconv.ParseFloat(arg, 64); err == nil {
		return arg
	}

	if _, err := strconv.ParseBool(arg); err == nil {
		return arg
	}

	return "`" + arg + "` is not a valid variable or literal"
}
