package prefix

import (
	"strings"

	. "github.com/spaceavocado/goillogical/internal"
	c "github.com/spaceavocado/goillogical/internal/expression/comparison"
)

func op(a string, b string) bool {
	return strings.HasPrefix(b, a)
}

func right(a string, b any) bool {
	switch b.(type) {
	case string:
		return op(a, b.(string))
	default:
		return false
	}
}

func left(a any, b any) bool {
	switch a.(type) {
	case string:
		return right(a.(string), b)
	default:
		return false
	}
}

func handler(evaluated []any) bool {
	return left(evaluated[0], evaluated[1])
}

func New(operator string, left Evaluable, right Evaluable) (Evaluable, error) {
	return c.New(operator, "<prefixes>", []Evaluable{left, right}, handler)
}
