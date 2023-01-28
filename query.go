package curry

import (
	"bytes"
	"curry/operators"
	"errors"
	"fmt"
	"strings"
)

var (
	InvalidNumberOfArguments = errors.New("an invalid number of arguments was generated - this is likely an internal error of the generator")
)

type Query struct {
	queryBase string
	where     *operators.Where
	tail      string
}

// Select begins the "select" query formation. The extra argument should contain all other query parts that should not pe parametrised, such as "join" statements - it's optional and can be left empty.
//
// All arguments are passed to the query as-is, so if anything requires special handling, like being surrounded with quotation marks, you need to do that yourself.
//
// Example: select * from users u join favourites f on u.id = f.user_id would be called Select("*", "users u", "join favourites f on u.id = f.user_id")
func Select(what string, from string, extra string) *Query {
	return &Query{queryBase: strings.TrimSpace(fmt.Sprintf("select %s from %s %s", what, from, extra))}
}

func (q *Query) Where(where *operators.Where) *Query {
	q.where = where
	return q
}

func (q *Query) Append(tail string) *Query {
	q.tail += tail
	return q
}

// ToExecutable returns the final query string and the list of arguments to use.
// The result of this method should be passed to your database connection handler.
func (q *Query) ToExecutable() (string, []interface{}, error) {
	arguments := make([]interface{}, 0)
	queryBuilder := &bytes.Buffer{}
	queryBuilder.WriteString(q.queryBase)

	if q.where != nil && q.where.ContainsAnything() {
		queryBuilder.WriteString(" where ")
		queryBuilder.WriteString(q.where.Print())
		arguments = append(arguments, q.where.GetOrderedArguments()...)
	}

	if q.tail != "" {
		queryBuilder.WriteString(" ")
		queryBuilder.WriteString(q.tail)
	}

	query := queryBuilder.String()
	parametrised, placeholders := replacePlaceholderArguments(query)
	if placeholders != len(arguments) {
		return "", nil, InvalidNumberOfArguments
	}

	return parametrised, arguments, nil
}

// replacePlaceholderArguments goes through the query and replaces each ? with $n, n being the next argument number starting from 1. The number of arguments (equal to the highest argument index) is then returned.
//
// If -1 is returned, then the query generation failed for an unknown reason. No expected errors are defined.
func replacePlaceholderArguments(query string) (string, int) {
	newQuery := &bytes.Buffer{}
	nextArgument := 1

	//this assumes we never want to just pass a ? to sql
	getIndex := func(q string) int { return strings.Index(q, "?") }

	for index := getIndex(query); index != -1; index = getIndex(query) {
		newQuery.WriteString(query[:index])
		if _, err := fmt.Fprintf(newQuery, "$%d", nextArgument); err != nil {
			return "", -1
		}
		nextArgument++
		query = query[(index + 1):]
	}

	newQuery.WriteString(query)
	return newQuery.String(), nextArgument - 1
}
