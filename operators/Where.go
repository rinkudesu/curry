package operators

import (
	"bytes"
	"fmt"
)

type WhereConnector int

const (
	And WhereConnector = iota
	Or
)

func (receiver *WhereConnector) Print() string {
	if *receiver == And {
		return "AND"
	}
	if *receiver == Or {
		return "OR"
	}
	return ""
}

type Where struct {
	properties         []SqlPrinter
	propertyConnectors []WhereConnector
}

var _ SqlPrinter = (*Where)(nil)

func WhereBegin(item *WhereItem) *Where {
	// connectors are always empty at the beginning
	where := &Where{propertyConnectors: make([]WhereConnector, 0)}
	if item.value == nil {
		where.properties = make([]SqlPrinter, 0)
	} else {
		where.properties = []SqlPrinter{item}
	}
	return where
}

func (w *Where) ContainsAnything() bool {
	return len(w.properties) > 0
}

func (w *Where) Or(printer SqlPrinter) *Where {
	if len(w.properties) > 0 {
		w.propertyConnectors = append(w.propertyConnectors, Or)
	}
	w.properties = append(w.properties, printer)
	return w
}

func (w *Where) And(printer SqlPrinter) *Where {
	if len(w.properties) > 0 {
		w.propertyConnectors = append(w.propertyConnectors, And)
	}
	w.properties = append(w.properties, printer)
	return w
}

func (w *Where) Print() string {
	// all arguments were optional and had the default value
	if len(w.properties) == 0 {
		return ""
	}

	query := &bytes.Buffer{}
	query.WriteString("(")

	for i, property := range w.properties {
		propertyString := property.Print()
		// ignore empty properties
		if propertyString == "" || propertyString == "()" {
			continue
		}

		if i > 0 {
			//ignore errors, as there's no sane ways of handling them anyway
			_, _ = fmt.Fprintf(query, " %s ", w.propertyConnectors[i-1].Print())
		}

		_, _ = fmt.Fprintf(query, "%s", propertyString)
	}

	query.WriteString(")")
	return query.String()
}

func (w *Where) GetOrderedArguments() []interface{} {
	interfaces := make([]interface{}, 0)
	for _, property := range w.properties {
		interfaces = append(interfaces, property.GetOrderedArguments()...)
	}
	return interfaces
}
