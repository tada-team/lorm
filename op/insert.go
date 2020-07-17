package op

import (
	"fmt"
	"log"
	"strings"
)

// [ WITH [ RECURSIVE ] with_query [, ...] ]
//INSERT INTO table_name [ AS alias ] [ ( column_name [, ...] ) ]
//    { DEFAULT VALUES | VALUES ( { expression | DEFAULT } [, ...] ) [, ...] | query }
//    [ ON CONFLICT [ conflict_target ] conflict_action ]
//    [ RETURNING * | output_expression [ [ AS ] output_name ] [, ...] ]
//
//where conflict_target can be one of:
//
//    ( { index_column_name | ( index_expression ) } [ COLLATE collation ] [ opclass ] [, ...] ) [ WHERE index_predicate ]
//    ON CONSTRAINT constraint_name
//
//and conflict_action is one of:
//
//    DO NOTHING
//    DO UPDATE SET { column_name = { expression | DEFAULT } |
//                    ( column_name [, ...] ) = ( { expression | DEFAULT } [, ...] ) |
//                    ( column_name [, ...] ) = ( sub-SELECT )
//                  } [, ...]
//              [ WHERE condition ]
type insertQuery struct {
	table      Table
	kv         Set
	onConflict string
	returning  string
}

func Insert(s Table, kv Set) insertQuery {
	return insertQuery{table: s, kv: kv}
}

func (q insertQuery) Returning(v ...Column) insertQuery {
	if len(v) > 0 {
		e := make([]Expr, len(v))
		for i := range v {
			e[i] = v[i].BareName()
		}
		q.returning = joinExpr(e, ",")
	}
	return q
}

func (q insertQuery) OnConflictDoNothing() insertQuery {
	q.onConflict = "ON CONFLICT DO NOTHING"
	return q
}

func (q insertQuery) String() string { return q.Query() }

func (q insertQuery) Query() string {
	if len(q.kv) == 0 {
		log.Panicln("invalid insertQuery:", q)
	}

	names := make([]string, 0)
	values := make([]string, 0)
	for _, v := range q.kv.SortedItems() {
		names = append(names, v.Column.BareName().String())
		values = append(values, v.Expr.String())
	}

	res := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", q.table, strings.Join(names, ", "), strings.Join(values, ", "))

	if q.onConflict != "" {
		res = fmt.Sprintf("%s %s", res, q.onConflict)
	}

	if q.returning != "" {
		res = fmt.Sprintf("%s RETURNING %s", res, q.returning)
	}

	return res
}
