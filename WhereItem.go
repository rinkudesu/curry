package curry

import "fmt"

type WhereItem struct {
	columnName string
	operator   string
	value      *Parameter
}

var _ SqlPrinter = (*WhereItem)(nil)

func (w *WhereItem) Print() string {
	if w.value == nil {
		return ""
	}
	return fmt.Sprintf("%s %s %s", w.columnName, w.operator, w.value.Print())
}

func (w *WhereItem) GetOrderedArguments() []interface{} {
	if w.value == nil {
		return []interface{}{}
	}
	return w.value.GetOrderedArguments()
}

func NewWhereItem(columnName string, operator string, value *Parameter) *WhereItem {
	return &WhereItem{columnName: columnName, operator: operator, value: value}
}
