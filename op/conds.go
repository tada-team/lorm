package op

type Conds []Expr

func (conds *Conds) Add(v Expr) {
	*conds = append(*conds, v)
}

func (conds Conds) Fork(v ...Expr) Conds {
	return append(v, conds...)
}

func (conds Conds) String() string {
	return And(conds...).String()
}
