package nil

import (
	e "github.com/spaceavocado/goillogical/evaluable"
	c "github.com/spaceavocado/goillogical/internal/expression/comparison"
)

func handler(evaluated []any) bool {
	return evaluated[0] == nil
}

func New(operator string, eval e.Evaluable) (e.Evaluable, error) {
	return c.New(operator, "<is nil>", []e.Evaluable{eval}, handler)
}
