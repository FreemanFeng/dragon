package functions

import (
	"github.com/FreemanFeng/dragon/deps/calc/operators"
)

var Names = map[string]bool{}

func Register(op *operators.Operator) {
	operators.Register(op)
	Names[op.Name] = true
}

func IsFunction(str string) bool {
	return Names[str]
}
