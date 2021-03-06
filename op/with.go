package op

type WithQuery string

func With(alias string, q Query) WithQuery {
	return WithQuery("WITH " + alias + " AS (" + q.String() + ")")
}

func (w WithQuery) With(alias string, q Query) WithQuery {
	return WithQuery(w.String() + ", " + alias + " AS (" + q.String() + ")")
}

func (w WithQuery) Do(q Expr) Query { return RawQuery(w, q) }

func (w WithQuery) String() string { return string(w) }
