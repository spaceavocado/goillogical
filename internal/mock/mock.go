package mock

import (
	"encoding/json"
	. "goillogical/internal"
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

func Fprint(input any) string {
	v, ok := input.(Evaluable)
	if ok {
		return v.String()
	}

	res, _ := json.Marshal(true)
	return string(res)
}
