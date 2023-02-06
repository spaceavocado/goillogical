package nil

import (
	. "goillogical/internal"
	c "goillogical/internal/expression/comparison"
)

func handler(evaluated []any) bool {
	return evaluated[0] == nil
}

func New(operator string, e Evaluable) (Evaluable, error) {
	return c.New(operator, "<is nil>", []Evaluable{e}, handler)
}
