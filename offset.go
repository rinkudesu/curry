package curry

type offset struct {
	value int
}

func (o *offset) Print() string {
	return "offset ?"
}

func (o *offset) GetOrderedArguments() []interface{} {
	return []interface{}{o.value}
}
