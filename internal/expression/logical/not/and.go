package and

import (
	"errors"
	. "goillogical/internal"
	l "goillogical/internal/expression/logical"
)

func handler(evaluated []bool) bool {
	for _, e := range evaluated {
		if e == false {
			return false
		}
	}
	return true
}

func New(operands []Evaluable) (Evaluable, error) {
	if len(operands) < 2 {
		return nil, errors.New("logical AND expression must have at least 2 operands")
	}

	return l.New("AND", operands, handler)
}
