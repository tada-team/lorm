package op

import "fmt"

type WithQuery string

func With(alias string, q Query) WithQuery {
	return WithQuery(fmt.Sprintf("WITH %s AS (%s)", alias, q))
}

func (w WithQuery) String() string { return string(w) }

func (w WithQuery) With(alias string, q Query) WithQuery {
	return WithQuery(fmt.Sprintf("%s, %s AS (%s)", w, alias, q))
}

func (w WithQuery) Do(q Expr) Query {
	return RawQuery(w, q)
}
