package op

import (
	"fmt"
	"log"
	"strings"
)

// [ WITH [ RECURSIVE ] with_query [, ...] ]
// UPDATE [ ONLY ] table_name [ * ] [ [ AS ] alias ]
//    SET { column_name = { expression | DEFAULT } |
//          ( column_name [, ...] ) = ( { expression | DEFAULT } [, ...] ) |
//          ( column_name [, ...] ) = ( sub-SELECT )
//        } [, ...]
//    [ FROM from_list ]
//    [ WHERE condition | WHERE CURRENT OF cursor_name ]
//    [ RETURNING * | output_expression [ [ AS ] output_name ] [, ...] ]
type updateQuery struct {
	table          Table
	kv             Set
	from           string
	whereCondition Expr
	returning      string
}

func Update(t Table, kv Set) updateQuery {
	return updateQuery{table: t, kv: kv}
}

func (q updateQuery) From(v ...string) updateQuery {
	q.from = strings.Join(v, ", ")
	return q
}

func (q updateQuery) Where(conds ...Expr) updateQuery {
	switch len(conds) {
	case 0:
		panic("empty where condition")
	case 1:
		q.whereCondition = conds[0]
	default:
		q.whereCondition = And(conds...)
	}
	return q
}

func (q updateQuery) Returning(values ...Expr) updateQuery {
	if len(values) == 0 {
		q.returning = "*"
	} else {
		for i, v := range values {
			if i > 0 {
				q.returning += ", "
			}
			if c, ok := v.(Column); ok {
				//q.returning += c.BareName().String()
				q.returning += c.String()
			} else {
				q.returning += v.String()
			}
		}
	}
	return q
}

func (q updateQuery) String() string { return q.Query() }

func (q updateQuery) Query() string {
	res := fmt.Sprintf("UPDATE %s SET %s", q.table, q.kv)
	if res == "" {
		log.Panicln("invalid updateQuery:", q)
	}
	if q.from != "" {
		res = fmt.Sprintf("%s FROM %s", res, q.from)
	}
	if q.whereCondition != nil && q.whereCondition.String() != "" {
		res = fmt.Sprintf("%s WHERE %s", res, q.whereCondition)
	}
	if q.returning != "" {
		res = fmt.Sprintf("%s RETURNING %s", res, q.returning)
	}
	return res
}
