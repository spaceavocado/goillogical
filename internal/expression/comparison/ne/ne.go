package ne

import (
	. "goillogical/internal"
	c "goillogical/internal/expression/comparison"
)

func handler(evaluated []any) bool {
	if !c.IsComparable(evaluated[0], evaluated[1]) {
		return true
	}
	return evaluated[0] != evaluated[1]
}

func New(left Evaluable, right Evaluable) (Evaluable, error) {
	return c.New("!=", []Evaluable{left, right}, handler)
}
