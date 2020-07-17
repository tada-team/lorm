package op

import (
	"fmt"
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

func (q deleteQuery) Query() string {
	res := fmt.Sprintf("DELETE FROM %s", q.table)
	if q.condition != nil && q.condition.String() != "" {
		res = fmt.Sprintf("%s WHERE %s", res, q.condition)
	}
	if q.outputExpression != nil && q.outputExpression.String() != "" {
		res = fmt.Sprintf("%s RETURNING %s", res, q.outputExpression)
	}
	return res
}
