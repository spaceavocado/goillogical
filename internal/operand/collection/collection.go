package collection

import (
	"errors"
	"fmt"

	e "github.com/spaceavocado/goillogical/evaluable"
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
	items []e.Evaluable
	opts  *SerializeOptions
}

func (c collection) Evaluate(ctx e.Context) (any, error) {
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

func (c collection) Serialize() any {
	head := c.items[0].Serialize()
	if shouldBeEscaped(head, c.opts) {
		head = escapeOperator(head.(string), c.opts)
	}
	res := []any{head}
	for i := 1; i < len(c.items); i++ {
		res = append(res, c.items[i].Serialize())
	}

	return res
}

func (c collection) Simplify(ctx e.Context) (any, e.Evaluable) {
	res := []any{}
	for _, i := range c.items {
		val, e := i.Simplify(ctx)
		if e != nil {
			return nil, &c
		}
		res = append(res, val)
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

func shouldBeEscaped(input any, opts *SerializeOptions) bool {
	if input == nil {
		return false
	}

	switch typed := input.(type) {
	case string:
		_, ok := opts.EscapedOperators[typed]
		return ok
	default:
		return false
	}
}

func escapeOperator(input string, opts *SerializeOptions) string {
	return fmt.Sprintf("%s%s", opts.EscapeCharacter, input)
}

func New(items []e.Evaluable, opts *SerializeOptions) (e.Evaluable, error) {
	if len(items) == 0 {
		return nil, errors.New("collection operand must have at least 1 item")
	}

	// for _, item := range items {
	// 	if isCollection(item) {
	// 		return nil, errors.New("collection cannot contain nested collection")
	// 	}
	// }

	return collection{items, opts}, nil
}
