package op

import (
	"strings"
	"sync/atomic"
)

// [ WITH [ RECURSIVE ] with_query [, ...] ]
// DELETE FROM [ ONLY ] table_name [ * ] [ [ AS ] alias ]
//    [ USING using_list ]
//    [ WHERE condition | WHERE CURRENT OF cursor_name ]
//    [ RETURNING * | output_expression [ [ AS ] output_name ] [, ...] ]
type deleteQuery struct {
	table            Table
	condition        Expr
	outputExpression Expr
}

func Delete(t Table) deleteQuery {
	return deleteQuery{table: t}
}

func (q deleteQuery) Where(conds ...Expr) deleteQuery {
	switch len(conds) {
	case 0:
		panic("empty where condition")
	case 1:
		q.condition = conds[0]
	default:
		q.condition = And(conds...)
	}
	return q
}

func (q deleteQuery) Returning(v Expr) deleteQuery {
	q.outputExpression = v
	return q
}

func (q deleteQuery) String() string { return q.Query() }

var deleteQueryMaxSize int32 = 16

func (q deleteQuery) Query() string {
	var b strings.Builder
	b.Grow(int(atomic.LoadInt32(&deleteQueryMaxSize)))

	b.WriteString("DELETE FROM ")
	b.WriteString(q.table.String())

	if q.condition != nil && q.condition.String() != "" {
		b.WriteString(" WHERE ")
		b.WriteString(q.condition.String())
	}
	if q.outputExpression != nil && q.outputExpression.String() != "" {
		b.WriteString(" RETURNING ")
		b.WriteString(q.outputExpression.String())
	}

	return maybeGrow(b.String(), &deleteQueryMaxSize)
}
