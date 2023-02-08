package eq

import (
	e "github.com/spaceavocado/goillogical/evaluable"
	c "github.com/spaceavocado/goillogical/internal/expression/comparison"
)

func handler(evaluated []any) bool {
	if !c.IsComparable(evaluated[0], evaluated[1]) {
		return false
	}
	return evaluated[0] == evaluated[1]
}

func New(operator string, left e.Evaluable, right e.Evaluable) (e.Evaluable, error) {
	return c.New(operator, "==", []e.Evaluable{left, right}, handler)
}
