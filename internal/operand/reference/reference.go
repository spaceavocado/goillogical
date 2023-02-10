package reference

import (
	"errors"
	"strconv"
	"strings"

	e "github.com/spaceavocado/goillogical/evaluable"

	"fmt"
	"regexp"
)

type SerializeOptions struct {
	// A function used to determine if the operand is a reference type, otherwise evaluated as a static value.
	//
	// - `true` = reference type
	// - `false` = value type
	From func(string) (string, error)
	// A function used to transform the operand into the reference annotation stripped form. I.e. remove any
	// annotation used to detect the reference type. E.g. "$Reference" => "Reference".
	To func(string) string
}

type SimplifyOptions struct {
	// Reference paths which should be ignored while simplification is applied. Must be an exact match.
	IgnoredPaths []string
	// Reference paths which should be ignored while simplification is applied. Matching regular expression patterns.
	IgnoredPathsRx []regexp.Regexp
}

func DefaultSerializeOptions() SerializeOptions {
	return SerializeOptions{
		From: func(path string) (string, error) {
			if len(path) > 1 && strings.HasPrefix(path, "$") {
				return path[1:], nil
			}
			return "", errors.New("invalid operand")
		},
		To: func(operand string) string {
			return fmt.Sprintf("$%s", operand)
		},
	}
}

type DataType string

const (
	Undefined DataType = "Undefined"
	Number    DataType = "Number"
	Integer   DataType = "Integer"
	Float     DataType = "Float"
	String    DataType = "String"
	Boolean   DataType = "Boolean"
)

type reference struct {
	addr    string
	path    string
	dt      DataType
	serOpts *SerializeOptions
	simOpts *SimplifyOptions
}

func (r reference) Evaluate(ctx e.Context) (any, error) {
	_, res, err := evaluate(ctx, r.path, r.dt)
	return res, err
}

func (r reference) Serialize() any {
	path := r.path

	if r.dt != Undefined {
		path = fmt.Sprintf("%s.(%s)", r.path, r.dt)
	}

	return r.serOpts.To(path)
}

func (r reference) Simplify(ctx e.Context) (any, e.Evaluable) {
	path, res, _ := evaluate(ctx, r.path, r.dt)
	if res != nil || isIgnoredPath(path, r.simOpts) {
		return res, nil
	}
	return nil, &r
}

func (r reference) String() string {
	return fmt.Sprintf("{%s}", r.addr)
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
			return Undefined, fmt.Errorf("unsupported \"%s\" type casting", matches[1])
		}
	}
	return Undefined, nil
}

func isIgnoredPath(path string, opts *SimplifyOptions) bool {
	for _, p := range opts.IgnoredPaths {
		if p == path {
			return true
		}
	}

	for _, r := range opts.IgnoredPathsRx {
		if r.MatchString(path) {
			return true
		}
	}

	return false
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

		return 0, fmt.Errorf("invalid conversion from from \"%s\" (string) to number", val)
	}

	switch typed := val.(type) {
	case int, float32, float64:
		return val, nil
	case string:
		return fromString(typed)
	case bool:
		if typed {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("invalid conversion from from \"%v\" to number", val)
	}
}

func toInteger(val any) (any, error) {
	switch typed := val.(type) {
	case int:
		return val, nil
	case float32:
		return int(typed), nil
	case float64:
		return int(typed), nil
	case string:
		res, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			return 0, fmt.Errorf("invalid conversion from from \"%v\" (string) to integer", val)
		}
		return int(res), nil
	case bool:
		if typed {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("invalid conversion from from \"%v\" to integer", val)
	}
}

func toFloat(val any) (any, error) {
	switch typed := val.(type) {
	case int:
		return float64(typed), nil
	case float32, float64:
		return val, nil
	case string:
		res, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			return 0, fmt.Errorf("invalid conversion from from \"%v\" (string) to float", val)
		}
		return res, nil
	default:
		return 0, fmt.Errorf("invalid conversion from from \"%v\" to float", val)
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
	switch typed := val.(type) {
	case int:
		if typed == 1 {
			return true, nil
		}
		if typed == 0 {
			return false, nil
		}
		return false, fmt.Errorf("invalid conversion from from \"%d\" to boolean", val)
	case string:
		term := strings.TrimSpace(strings.ToLower(val.(string)))
		if term == "true" || term == "1" {
			return true, nil
		}
		if term == "false" || term == "0" {
			return false, nil
		}
		return false, fmt.Errorf("invalid conversion from from \"%s\" to boolean", val)
	case bool:
		return val.(bool), nil
	default:
		return false, fmt.Errorf("invalid conversion from from \"%v\" to boolean", val)
	}
}

func contextLookup(ctx e.Context, path string) (string, any) {
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

func evaluate(ctx e.Context, path string, dt DataType) (string, any, error) {
	resolvedPath, value := contextLookup(ctx, path)

	if value == nil {
		return resolvedPath, nil, nil
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

func New(addr string, serOpts *SerializeOptions, simOpts *SimplifyOptions) (e.Evaluable, error) {
	dt, err := getDataType(addr)
	if err != nil {
		return nil, err
	}

	return reference{addr, trimDataType(addr), dt, serOpts, simOpts}, nil
}
