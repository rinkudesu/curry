package curry

type limit struct {
	value int
}

func (l *limit) Print() string {
	return "limit ?"
}

func (l *limit) GetOrderedArguments() []interface{} {
	return []interface{}{l.value}
}
