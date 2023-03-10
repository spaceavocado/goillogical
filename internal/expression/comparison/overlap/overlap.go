package overlap

import (
	"reflect"

	e "github.com/spaceavocado/goillogical/evaluable"
	c "github.com/spaceavocado/goillogical/internal/expression/comparison"
)

func handler(evaluated []any) bool {
	t1s := reflect.TypeOf(evaluated[0]).Kind() == reflect.Slice
	t2s := reflect.TypeOf(evaluated[1]).Kind() == reflect.Slice

	if !t1s || !t2s {
		return false
	}

	left := reflect.ValueOf(evaluated[0])
	right := reflect.ValueOf(evaluated[1])

	for i := 0; i < left.Len(); i++ {
		p1 := left.Index(i).Interface()
		for j := 0; j < right.Len(); j++ {
			p2 := right.Index(j).Interface()
			if c.IsComparable(p1, p2) && p1 == p2 {
				return true
			}
		}
	}

	return false
}

func New(operator string, left e.Evaluable, right e.Evaluable) (e.Evaluable, error) {
	return c.New(operator, "<overlaps>", []e.Evaluable{left, right}, handler)
}
