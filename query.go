package curry

type Query struct {
	queryBase string
}

func Base(queryBase string) *Query {
	return &Query{queryBase: queryBase}
}

func (q *Query) ToString() string {
	return q.queryBase
}
