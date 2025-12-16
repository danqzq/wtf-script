package types

import (
	"wtf-script/config"
)

type IInterpreter interface {
	GetConfig() *config.Config
	GenerateRandomString(n int, charset string) string
	SetSeed(seed int64)
}

type IBuiltinFunc func(args []any, i IInterpreter) any
