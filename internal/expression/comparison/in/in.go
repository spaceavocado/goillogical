package in

import (
	"reflect"

	e "github.com/spaceavocado/goillogical/evaluable"
	c "github.com/spaceavocado/goillogical/internal/expression/comparison"
)

func handler(evaluated []any) bool {
	t1s := reflect.TypeOf(evaluated[0]).Kind() == reflect.Slice
	t2s := reflect.TypeOf(evaluated[1]).Kind() == reflect.Slice

	if (t1s && t2s) || (!t1s && !t2s) {
		return false
	}

	var haystack reflect.Value
	var needle any
	if t1s {
		haystack = reflect.ValueOf(evaluated[0])
		needle = evaluated[1]
	} else {
		haystack = reflect.ValueOf(evaluated[1])
		needle = evaluated[0]
	}

	for i := 0; i < haystack.Len(); i++ {
		p := haystack.Index(i).Interface()
		if c.IsComparable(p, needle) && haystack.Index(i).Interface() == needle {
			return true
		}
	}

	return false
}

func New(operator string, left e.Evaluable, right e.Evaluable) (e.Evaluable, error) {
	return c.New(operator, "<in>", []e.Evaluable{left, right}, handler)
}
