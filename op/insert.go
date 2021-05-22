package op

import (
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
	kvs        []Set
	onConflict string
	returning  string
}

func Insert(s Table, kvs ...Set) insertQuery {
	return insertQuery{table: s, kvs: kvs}
}

func (q insertQuery) Returning(v ...Column) insertQuery {
	if len(v) > 0 {
		var b strings.Builder
		for i, c := range v {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(c.BareName().String())
		}
		q.returning = b.String()
	}
	return q
}

func (q insertQuery) OnConflictDoNothing() insertQuery {
	q.onConflict = "ON CONFLICT DO NOTHING"
	return q
}

func (q insertQuery) String() string { return q.Query() }

var insertQueryMaxSize = 200

func (q insertQuery) Query() string {
	if len(q.kvs) == 0 {
		log.Panicln("invalid insertQuery:", q)
	}

	//items := q.kv.SortedItems()

	var b strings.Builder
	b.Grow(insertQueryMaxSize)

	b.WriteString("INSERT INTO ")
	b.WriteString(q.table.String())

	b.WriteString(" (")
	for i, item := range q.kvs[0].SortedItems() {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(item.Column.BareName().String())
	}
	b.WriteString(") VALUES ")
	for j, kv := range q.kvs {
		if j > 0 {
			b.WriteString(", ")
		}
		b.WriteString("(")
		items := kv.SortedItems()
		for i, item := range items {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(item.Expr.String())
		}
		b.WriteString(")")
	}

	if q.onConflict != "" {
		b.WriteString(" ")
		b.WriteString(q.onConflict)
	}

	if q.returning != "" {
		b.WriteString(" RETURNING ")
		b.WriteString(q.returning)
	}

	return maybeGrow(b.String(), &insertQueryMaxSize)
}
