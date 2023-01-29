package curry

type SqlPrinter interface {
	Print() string
	GetOrderedArguments() []interface{}
}
