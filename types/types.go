package types

import (
	"wtf-script/config"
)

type IInterpreter interface {
	GetVariable(name string) (*Variable, bool)
	SetVariable(name string, value any)
	SetSeed(seed int64)
	GetConfig() *config.Config
	GenerateRandomString(n int, charset string) string
}

type IBuiltinFunc func(args []any, i IInterpreter) any
