package eq

import (
	. "github.com/spaceavocado/goillogical/internal"
	c "github.com/spaceavocado/goillogical/internal/expression/comparison"
)

func handler(evaluated []any) bool {
	if !c.IsComparable(evaluated[0], evaluated[1]) {
		return false
	}
	return evaluated[0] == evaluated[1]
}

func New(operator string, left Evaluable, right Evaluable) (Evaluable, error) {
	return c.New(operator, "==", []Evaluable{left, right}, handler)
}
