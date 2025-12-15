package interpreter

import (
	"math/rand"
	"strings"
	"time"
	"wtf-script/builtins"
	"wtf-script/types"
)

type Interpreter struct {
	Variables map[string]types.Variable
	Builtins  map[string]func(args []string, i types.IInterpreter)
	Rand      *rand.Rand
}

func (i *Interpreter) GetVariable(name string) (*types.Variable, bool) {
	v, ok := i.Variables[name]
	return &v, ok
}

func (i *Interpreter) SetVariable(name string, value interface{}) {
	if v, ok := i.Variables[name]; ok {
		v.Value = value
		i.Variables[name] = v
	}
}

func (i *Interpreter) SetSeed(seed int64) {
	i.Rand.Seed(seed)
}

func NewInterpreter() *Interpreter {
	i := &Interpreter{
		Variables: make(map[string]types.Variable),
		Builtins:  make(map[string]func(args []string, i types.IInterpreter)),
		Rand:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	builtins.RegisterBuiltins(func(name string, fn types.IBuiltinFunc) {
		i.Builtins[name] = fn
	})
	return i
}

func (i *Interpreter) Run(code string) {
	lines := strings.Split(code, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			//i.parseLine(line)
		}
	}
}
