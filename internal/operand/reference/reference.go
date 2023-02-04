package reference

import (
	"errors"
	. "goillogical/internal"
	"strconv"
	"strings"

	"fmt"
	"reflect"
	"regexp"
)

type DataType byte

const (
	Unknown DataType = iota
	Number
	Integer
	Float
	String
	Boolean
)

type reference struct {
	addr string
	path string
	dt   DataType
}

// func (v reference) Kind() Kind {
// 	return Reference
// }

func (r reference) String() string {
	return fmt.Sprintf("{%s}", r.addr)
}

func (r reference) Evaluate(ctx Context) (any, error) {
	_, res, err := evaluate(ctx, r.path, r.dt)
	return res, err
}

func getDataType(path string) (DataType, error) {
	re := regexp.MustCompile(`^.+\.\(([A-Z][a-z]+)\)$`)
	matches := re.FindStringSubmatch(path)
	if len(matches) > 1 {
		switch matches[1] {
		case "Number":
			return Number, nil
		case "Integer":
			return Integer, nil
		case "Float":
			return Float, nil
		case "String":
			return String, nil
		case "Boolean":
			return Boolean, nil
		default:
			return Unknown, errors.New(fmt.Sprintf("unsupported \"%s\" type casting", matches[1]))
		}
	}
	return Unknown, nil
}

func trimDataType(path string) string {
	re := regexp.MustCompile(`.\(([A-Z][a-z]+)\)$`)
	return re.ReplaceAllString(path, "")
}

func toNumber(val any) (any, error) {
	reFloat := regexp.MustCompile(`^\d+\.\d+$`)
	reInt := regexp.MustCompile(`^0$|^[1-9]\d*$`)

	fromString := func(val string) (any, error) {
		if reFloat.MatchString(val) {
			result, _ := strconv.ParseFloat(val, 64)
			return result, nil
		}
		if reInt.MatchString(val) {
			result, _ := strconv.Atoi(val)
			return result, nil
		}

		return 0, errors.New(fmt.Sprintf("invalid conversion from from \"%s\" (string) to number", val))
	}

	switch val.(type) {
	case int, float32, float64:
		return val, nil
	case string:
		return fromString(val.(string))
	case bool:
		if val.(bool) == true {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, errors.New(fmt.Sprintf("invalid conversion from from \"%v\" to number", val))
	}
}

func toInteger(val any) (any, error) {
	switch val.(type) {
	case int:
		return val, nil
	case float32:
		return int(val.(float32)), nil
	case float64:
		return int(val.(float64)), nil
	case string:
		res, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			return 0, errors.New(fmt.Sprintf("invalid conversion from from \"%v\" (string) to integer", val))
		}
		return int(res), nil
	case bool:
		if val.(bool) == true {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, errors.New(fmt.Sprintf("invalid conversion from from \"%v\" to integer", val))
	}
}

func toFloat(val any) (any, error) {
	switch val.(type) {
	case int:
		return float64(val.(int)), nil
	case float32, float64:
		return val, nil
	case string:
		res, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			return 0, errors.New(fmt.Sprintf("invalid conversion from from \"%v\" (string) to float", val))
		}
		return res, nil
	default:
		return 0, errors.New(fmt.Sprintf("invalid conversion from from \"%v\" to float", val))
	}
}

func toString(val any) string {
	switch val.(type) {
	case int:
		return fmt.Sprintf("%d", val)
	case float32, float64:
		return fmt.Sprintf("%f", val)
	case string:
		return fmt.Sprintf("%s", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func toBoolean(val any) (bool, error) {
	switch val.(type) {
	case int:
		if val.(int) == 1 {
			return true, nil
		}
		if val.(int) == 0 {
			return false, nil
		}
		return false, errors.New(fmt.Sprintf("invalid conversion from from \"%d\" to boolean", val))
	case string:
		term := strings.TrimSpace(strings.ToLower(val.(string)))
		if term == "true" || term == "1" {
			return true, nil
		}
		if term == "false" || term == "0" {
			return false, nil
		}
		return false, errors.New(fmt.Sprintf("invalid conversion from from \"%s\" to boolean", val))
	case bool:
		return val.(bool), nil
	default:
		return false, errors.New(fmt.Sprintf("invalid conversion from from \"%v\" to boolean", val))
	}
}

func flattenContext(ctx Context) map[string]any {
	res := make(map[string]any)
	var lookup func(p any, path string)

	joinPath := func(a string, b string) string {
		if len(a) == 0 {
			return b
		}
		return fmt.Sprintf("%s.%s", a, b)
	}

	lookup = func(val any, path string) {
		v := reflect.ValueOf(val)
		switch v.Kind() {
		case reflect.Bool:
			fallthrough
		case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
			fallthrough
		case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
			fallthrough
		case reflect.Float32, reflect.Float64:
			fallthrough
		case reflect.String:
			res[path] = val
			break
		case reflect.Map:
			for prop, val := range val.(map[string]any) {
				lookup(val, joinPath(path, prop))
			}
			break
		case reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				lookup(v.Index(i).Interface(), fmt.Sprintf("%s[%d]", path, i))
			}
			break
		default:
			return
		}
	}

	lookup(ctx, "")
	return res
}

func contextLookup(ctx Context, path string) (string, any) {
	rxPath := regexp.MustCompile(`{([^{}]+)}`)
	for match := rxPath.FindStringSubmatchIndex(path); len(match) > 0; {
		_, val := contextLookup(ctx, string(path[match[2]:match[3]]))
		if val == nil {
			return path, nil
		}
		path = path[0:match[0]] + fmt.Sprintf("%v", val) + path[match[1]:]
		match = rxPath.FindStringSubmatchIndex(path)
	}

	if val, ok := ctx[path]; ok {
		return path, val
	}

	return path, nil
}

func isEvaluatedValue(value any) bool {
	if value == nil {
		return false
	}
	switch value.(type) {
	case int, float32, float64, bool, string:
		return true
	default:
		return false
	}
}

func evaluate(ctx Context, path string, dt DataType) (string, any, error) {
	resolvedPath, value := contextLookup(flattenContext(ctx), path)

	if !isEvaluatedValue(value) {
		return resolvedPath, nil, errors.New(fmt.Sprintf("invalid evaluated value at \"%s\" path", path))
	}

	switch dt {
	case Number:
		val, err := toNumber(value)
		return resolvedPath, val, err
	case Integer:
		val, err := toInteger(value)
		return resolvedPath, val, err
	case Float:
		val, err := toFloat(value)
		return resolvedPath, val, err
	case Boolean:
		val, err := toBoolean(value)
		return resolvedPath, val, err
	case String:
		return resolvedPath, toString(value), nil
	default:
		return resolvedPath, value, nil
	}
}

func New(addr string) (Evaluable, error) {
	dt, err := getDataType(addr)
	if err != nil {
		return nil, err
	}

	return reference{addr: addr, path: trimDataType(addr), dt: dt}, nil
}
