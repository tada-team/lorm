package op

type TsVector string

type TsQuery string

func (v TsVector) Search(q TsQuery) Expr {
	return Raw(v, "@@", q)
}
