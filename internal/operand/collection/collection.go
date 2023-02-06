package collection

import (
	"errors"
	. "goillogical/internal"
)

type SerializeOptions struct {
	EscapedOperators map[string]bool
	EscapeCharacter  string
}

func DefaultSerializeOptions() SerializeOptions {
	return SerializeOptions{
		EscapedOperators: map[string]bool{},
		EscapeCharacter:  "\\",
	}
}

type collection struct {
	items []Evaluable
}

func (c collection) Kind() Kind {
	return Collection
}

func (c collection) Evaluate(ctx Context) (any, error) {
	res := make([]any, len(c.items))
	for i, item := range c.items {
		val, err := item.Evaluate((ctx))
		if err != nil {
			return nil, err
		}
		res[i] = val
	}
	return res, nil
}

func (c collection) String() string {
	res := "["
	for i, item := range c.items {
		res += item.String()
		if i < len(c.items)-1 {
			res += ", "
		}
	}
	res += "]"
	return res
}

func New(items []Evaluable) (Evaluable, error) {
	if len(items) == 0 {
		return nil, errors.New("collection operand must have at least 1 item")
	}

	// for _, item := range items {
	// 	if isCollection(item) {
	// 		return nil, errors.New("collection cannot contain nested collection")
	// 	}
	// }

	return collection{items}, nil
}
