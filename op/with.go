package op

import "fmt"

type with string

func With(alias string, q Query) with {
	return with(fmt.Sprintf("WITH %s AS (%s)", alias, q))
}

func (w with) String() string { return string(w) }

func (w with) With(alias string, q Query) with {
	return with(fmt.Sprintf("%s, %s AS (%s)", w, alias, q))
}

func (w with) Do(q Expr) Query {
	return RawQuery(w, q)
}
