package lorm

import (
	"github.com/tada-team/lorm/op"
)

type Filter interface {
	Transactional
	GetArgs() op.Args
	GetConds() op.Conds
	IsEmpty() bool
	GetLock() op.Lock
	GetOrderBy() op.Expr
	GetLimit() int
	GetOffset() int
}

type BaseFilter struct {
	BaseTransactional
	Conds   op.Conds
	Args    op.Args
	empty   bool
	limit   int
	offset  int
	lock    op.Lock
	orderBy op.Expr
}

func (f BaseFilter) GetArgs() op.Args    { return f.Args }
func (f BaseFilter) GetConds() op.Conds  { return f.Conds }
func (f BaseFilter) IsEmpty() bool       { return f.empty }
func (f BaseFilter) GetLock() op.Lock    { return f.lock }
func (f BaseFilter) GetOrderBy() op.Expr { return f.orderBy }
func (f BaseFilter) GetLimit() int       { return f.limit }
func (f BaseFilter) GetOffset() int      { return f.offset }

func (f *BaseFilter) SetOrderBy(v op.Expr)         { f.orderBy = v }
func (f *BaseFilter) SetLimit(v int)               { f.limit = v }
func (f *BaseFilter) SetOffset(v int)              { f.offset = v }
func (f *BaseFilter) SetEmpty()                    { f.empty = true }
func (f *BaseFilter) SetLock(tx *Tx, lock op.Lock) { f.SetTx(tx); f.lock = lock }

// shortcuts

// C = condition
func (f *BaseFilter) C(v op.Expr) { f.Conds.Add(v) }

// EC = empty condition
func (f *BaseFilter) EC(empty bool, v op.Expr) {
	if empty {
		f.SetEmpty()
	} else {
		f.Conds.Add(v)
	}
}

// MC = maybe condition
func (f *BaseFilter) MC(apply bool, v op.Expr) {
	if apply {
		f.Conds.Add(v)
	}
}
