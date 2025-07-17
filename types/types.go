package types

type IInterpreter interface {
	GetVariable(name string) (*Variable, bool)
	SetVariable(name string, value interface{})
	SetSeed(seed int64)
}

type IBuiltinFunc func(args []string, i IInterpreter)
