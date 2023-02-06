package internal

type evaluableMock struct {
	val any
	str string
}

func (m evaluableMock) Kind() Kind {
	return Unknown
}

func (m evaluableMock) String() string {
	return m.str
}

func (m evaluableMock) Evaluate(ctx Context) (any, error) {
	return m.val, nil
}

func (m evaluableMock) Serialize() any {
	return m.val
}

func EvaluableMock(val any, str string) Evaluable {
	return evaluableMock{val, str}
}
