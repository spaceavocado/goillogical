package test

import (
	"encoding/json"
	. "goillogical/internal"
)

func Fprint(input any) string {
	v, ok := input.(Evaluable)
	if ok {
		return v.String()
	}

	res, _ := json.Marshal(input)
	return string(res)
}
