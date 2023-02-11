package curry

import (
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
	query := Select("*", "users", "").Where(WhereBegin(NewWhereItem("id", "=", NewParameter(1))))

	result, parameters, err := query.ToExecutable()

	expectedResult := "select * from users where (id = $1)"
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, 1, len(parameters))
	assert.Equal(t, 1, parameters[0])
}

func TestQuery_ToExecutable_SelectWithEmptyWhere(t *testing.T) {
	t.Parallel()
	query := Select("*", "users", "").Where(WhereBegin(NewWhereItem("id", "=", NewOptionalParameter[int](1, 1))))

	result, parameters, err := query.ToExecutable()

	expectedResult := "select * from users"
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Empty(t, parameters)
}

func TestQuery_ToExecutable_SelectWithWhereAndOr(t *testing.T) {
	t.Parallel()
	query := Select("*", "users", "").Where(WhereBegin(NewWhereItem("id", "=", NewParameter(1))).And(NewWhereItem("name", "=", NewParameter("test"))).Or(NewWhereItem("test", "!=", NewParameter("000"))))

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
	query := Select("*", "users", "").Where(WhereBegin(NewWhereItem("id", "=", NewOptionalParameter[int](1, 1))).And(NewWhereItem("test", "=", NewParameter("a"))).Or(NewWhereItem("aaa", "=", NewOptionalParameter(1, 1))))

	result, parameters, err := query.ToExecutable()

	expectedResult := "select * from users where (test = $1)"
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, 1, len(parameters))
	assert.Equal(t, "a", parameters[0])
}

func TestQuery_ToExecutable_SelectWithNestedWhere(t *testing.T) {
	t.Parallel()
	query := Select("*", "users", "").Where(WhereBegin(NewWhereItem("id", "=", NewParameter(1))).And(WhereBegin(NewWhereItem("test", "=", NewParameter(15))).Or(NewWhereItem("test", "=", NewParameter(2)))))

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
	query := Select("*", "users", "").Where(WhereBegin(NewWhereItem("id", "=", NewParameter(1))).And(WhereBegin(NewWhereItem("test", "=", NewOptionalParameter(1, 1)))))

	result, parameters, err := query.ToExecutable()

	expectedResult := "select * from users where (id = $1)"
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, 1, len(parameters))
	assert.Equal(t, 1, parameters[0])
}

func TestQuery_ToExecutable_SelectWithSimpleSingleWhereLimitOffsetAppend(t *testing.T) {
	t.Parallel()
	query := Select("*", "users", "").Where(WhereBegin(NewWhereItem("id", "=", NewParameter(1)))).Limit(2).Offset(3).Append("returning id")

	result, parameters, err := query.ToExecutable()

	expectedResult := "select * from users where (id = $1) offset $2 limit $3 returning id"
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, 3, len(parameters))
	assert.Equal(t, 1, parameters[0])
	assert.Equal(t, 3, parameters[1])
	assert.Equal(t, 2, parameters[2])
}

func TestQuery_ToExecutable_LimitOffsetZerosIgnored(t *testing.T) {
	t.Parallel()
	query := Select("*", "users", "").Where(WhereBegin(NewWhereItem("id", "=", NewParameter(1)))).Limit(0).Offset(0).Append("returning id")

	result, parameters, err := query.ToExecutable()

	expectedResult := "select * from users where (id = $1) returning id"
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, 1, len(parameters))
	assert.Equal(t, 1, parameters[0])
}

func TestQuery_ToExecutable_SelectWithSimpleSingleWhereLimitOffsetAppendOrderBy(t *testing.T) {
	t.Parallel()
	query := Select("*", "users", "").Where(WhereBegin(NewWhereItem("id", "=", NewParameter(1)))).Limit(2).Offset(3).Append("returning id").OrderBy("test")

	result, parameters, err := query.ToExecutable()

	expectedResult := "select * from users where (id = $1) order by test offset $2 limit $3 returning id"
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, 3, len(parameters))
	assert.Equal(t, 1, parameters[0])
	assert.Equal(t, 3, parameters[1])
	assert.Equal(t, 2, parameters[2])
}

func TestQuery_ToExecutable_SelectWithSimpleSingleWhereLimitOffsetAppendOrderByDesc(t *testing.T) {
	t.Parallel()
	query := Select("*", "users", "").Where(WhereBegin(NewWhereItem("id", "=", NewParameter(1)))).Limit(2).Offset(3).Append("returning id").OrderByDescending("test")

	result, parameters, err := query.ToExecutable()

	expectedResult := "select * from users where (id = $1) order by test desc offset $2 limit $3 returning id"
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, 3, len(parameters))
	assert.Equal(t, 1, parameters[0])
	assert.Equal(t, 3, parameters[1])
	assert.Equal(t, 2, parameters[2])
}
