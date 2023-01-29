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

var tests = []testData{}

func TestQuery_ToExecutable_SimpleSelect(t *testing.T) {
	t.Parallel()
	query := Select("*", "users", "")

	result, parameters, err := query.ToExecutable()

	expectedResult := "select * from users"
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Empty(t, parameters)
}

func TestQuery_ToExecutable_SelectWithJoin(t *testing.T) {
	t.Parallel()
	query := Select("*", "users u", "join favourites f on f.user_id = u.id")

	result, parameters, err := query.ToExecutable()

	expectedResult := "select * from users u join favourites f on f.user_id = u.id"
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Empty(t, parameters)
}

func TestQuery_ToExecutable_SelectWithSimpleSingleWhere(t *testing.T) {
	t.Parallel()
	query := Select("*", "users", "").Where(operators.WhereBegin(operators.NewWhereItem("id", "=", operators.NewParameter(1))))

	result, parameters, err := query.ToExecutable()

	expectedResult := "select * from users where (id = $1)"
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, 1, len(parameters))
	assert.Equal(t, 1, parameters[0])
}

func TestQuery_ToExecutable_SelectWithEmptyWhere(t *testing.T) {
	t.Parallel()
	query := Select("*", "users", "").Where(operators.WhereBegin(operators.NewWhereItem("id", "=", operators.NewOptionalParameter[int](1, 1))))

	result, parameters, err := query.ToExecutable()

	expectedResult := "select * from users"
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Empty(t, parameters)
}

func TestQuery_ToExecutable_SelectWithWhereAndOr(t *testing.T) {
	t.Parallel()
	query := Select("*", "users", "").Where(operators.WhereBegin(operators.NewWhereItem("id", "=", operators.NewParameter(1))).And(operators.NewWhereItem("name", "=", operators.NewParameter("test"))).Or(operators.NewWhereItem("test", "!=", operators.NewParameter("000"))))

	result, parameters, err := query.ToExecutable()

	expectedResult := "select * from users where (id = $1 AND name = $2 OR test != $3)"
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, 3, len(parameters))
	assert.Equal(t, 1, parameters[0])
	assert.Equal(t, "test", parameters[1])
	assert.Equal(t, "000", parameters[2])
}

func TestQuery_ToExecutable_SelectWithWhereIncludingEmpty(t *testing.T) {
	t.Parallel()
	query := Select("*", "users", "").Where(operators.WhereBegin(operators.NewWhereItem("id", "=", operators.NewOptionalParameter[int](1, 1))).And(operators.NewWhereItem("test", "=", operators.NewParameter("a"))).Or(operators.NewWhereItem("aaa", "=", operators.NewOptionalParameter(1, 1))))

	result, parameters, err := query.ToExecutable()

	expectedResult := "select * from users where (test = $1)"
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, 1, len(parameters))
	assert.Equal(t, "a", parameters[0])
}

func TestQuery_ToExecutable_SelectWithNestedWhere(t *testing.T) {
	t.Parallel()
	query := Select("*", "users", "").Where(operators.WhereBegin(operators.NewWhereItem("id", "=", operators.NewParameter(1))).And(operators.WhereBegin(operators.NewWhereItem("test", "=", operators.NewParameter(15))).Or(operators.NewWhereItem("test", "=", operators.NewParameter(2)))))

	result, parameters, err := query.ToExecutable()

	expectedResult := "select * from users where (id = $1 AND (test = $2 OR test = $3))"
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, 3, len(parameters))
	assert.Equal(t, 1, parameters[0])
	assert.Equal(t, 15, parameters[1])
	assert.Equal(t, 2, parameters[2])
}

func TestQuery_ToExecutable_SelectWithEmptyNestedWhere(t *testing.T) {
	t.Parallel()
	query := Select("*", "users", "").Where(operators.WhereBegin(operators.NewWhereItem("id", "=", operators.NewParameter(1))).And(operators.WhereBegin(operators.NewWhereItem("test", "=", operators.NewOptionalParameter(1, 1)))))

	result, parameters, err := query.ToExecutable()

	expectedResult := "select * from users where (id = $1)"
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, 1, len(parameters))
	assert.Equal(t, 1, parameters[0])
}
