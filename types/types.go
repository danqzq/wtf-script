package types

import (
	"wtf-script/config"
)

type IInterpreter interface {
	GetConfig() *config.Config
	GenerateRandomString(n int, charset string) string
	SetSeed(seed int64)
	LogError(format string, args ...any)
}

type IBuiltinFunc func(args []any, i IInterpreter) any
