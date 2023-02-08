package test

import (
	"encoding/json"

	e "github.com/spaceavocado/goillogical/evaluable"
)

func Fprint(input any) string {
	v, ok := input.(e.Evaluable)
	if ok {
		return v.String()
	}

	res, _ := json.Marshal(input)
	return string(res)
}
