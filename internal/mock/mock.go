package mock

import (
	. "goillogical/internal"
	collection "goillogical/internal/operand/collection"
	reference "goillogical/internal/operand/reference"
	value "goillogical/internal/operand/value"
	"regexp"
)

type eMock struct {
	val any
	str string
}

func (m eMock) Kind() Kind {
	return Unknown
}

func (m eMock) String() string {
	return m.str
}

func (m eMock) Evaluate(ctx Context) (any, error) {
	return m.val, nil
}

func (m eMock) Serialize() any {
	return m.val
}

func (m eMock) Simplify(ctx Context) (any, Evaluable) {
	return m.val, nil
}

func E(val any, str string) Evaluable {
	return eMock{val, str}
}

func Val(val any) Evaluable {
	e, _ := value.New(val)
	return e
}

func Ref(val string) Evaluable {
	serOpts := reference.DefaultSerializeOptions()
	simOpts := reference.SimplifyOptions{
		IgnoredPaths:   []string{"ignored"},
		IgnoredPathsRx: []regexp.Regexp{},
	}
	e, _ := reference.New(val, &serOpts, &simOpts)
	return e
}

func Col(items ...Evaluable) Evaluable {
	opts := collection.DefaultSerializeOptions()
	e, _ := collection.New(items, &opts)
	return e
}
