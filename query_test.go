package curry

import (
	"curry/operators"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testData struct {
	name              string
	query             *Query
	expectedResult    string
	expectedArguments []interface{}
}

var tests = []testData{
	{name: "simple select", query: Select("*", "users", ""), expectedResult: "select * from users"},
	{name: "select with join", query: Select("*", "users u", "join favourites f on f.user_id = u.id"), expectedResult: "select * from users u join favourites f on f.user_id = u.id"},
	{name: "select with simple where (single)", query: Select("*", "users", "").Where(operators.WhereBegin(operators.NewWhereItem("id", "=", operators.NewParameter(1)))), expectedResult: "select * from users where (id = $1)"},
	{name: "select with empty where", query: Select("*", "users", "").Where(operators.WhereBegin(operators.NewWhereItem("id", "=", operators.NewOptionalParameter[int](1, 1)))), expectedResult: "select * from users"},
	{name: "select with simple where (and or)", query: Select("*", "users", "").Where(operators.WhereBegin(operators.NewWhereItem("id", "=", operators.NewParameter(1))).And(operators.NewWhereItem("name", "=", operators.NewParameter("test"))).Or(operators.NewWhereItem("test", "!=", operators.NewParameter("000")))), expectedResult: "select * from users where (id = $1 AND name = $2 OR test != $3)"},
	{name: "select with where including empty", query: Select("*", "users", "").Where(operators.WhereBegin(operators.NewWhereItem("id", "=", operators.NewOptionalParameter[int](1, 1))).And(operators.NewWhereItem("test", "=", operators.NewParameter("a"))).Or(operators.NewWhereItem("aaa", "=", operators.NewOptionalParameter(1, 1)))), expectedResult: "select * from users where (test = $1)"},
	{name: "select with nested where", query: Select("*", "users", "").Where(operators.WhereBegin(operators.NewWhereItem("id", "=", operators.NewParameter(1))).And(operators.WhereBegin(operators.NewWhereItem("test", "=", operators.NewParameter(1))).Or(operators.NewWhereItem("test", "=", operators.NewParameter(2))))), expectedResult: "select * from users where (id = $1 AND (test = $2 OR test = $3))"},
	{name: "select with empty nested where", query: Select("*", "users", "").Where(operators.WhereBegin(operators.NewWhereItem("id", "=", operators.NewParameter(1))).And(operators.WhereBegin(operators.NewWhereItem("test", "=", operators.NewOptionalParameter(1, 1))))), expectedResult: "select * from users where (id = $1)"},
}

// I'm gonna break everything you know about unit tests here and just make a single test method comparing a bunch of function calls to raw sql queries, as it's much easier to parse mentally
func TestQuery_ToExecutable_QueryCompare(t *testing.T) {
	for i := range tests {
		index := i
		t.Run(tests[index].name, func(t *testing.T) {
			t.Parallel()
			result, _, err := tests[index].query.ToExecutable()
			assert.Nil(t, err)
			if tests[index].expectedArguments != nil {
				//todo
			}
			assert.Equal(t, tests[index].expectedResult, result)
		})
	}
}
