package curry

type Parameter struct {
	value interface{}
}

func NewParameter(value interface{}) *Parameter {
	return &Parameter{value: value}
}

var _ SqlPrinter = (*Parameter)(nil)

// NewOptionalParameter returns the parameter if value not equal to ifNot. If they are equal, nil is returned instead.
func NewOptionalParameter[T comparable](value T, ifNot T) *Parameter {
	if value == ifNot {
		return nil
	}

	return &Parameter{value: value}
}

// Print in this case will always return a question mark, as parameters are not to be placed directly in sql code
func (p *Parameter) Print() string {
	return "?"
}

func (p *Parameter) GetOrderedArguments() []interface{} {
	return []interface{}{p.value}
}
