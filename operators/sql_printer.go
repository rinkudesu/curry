package operators

type SqlPrinter interface {
	Print() string
	GetOrderedArguments() []interface{}
}
