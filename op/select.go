package op

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Lock string

var selectQueryMaxSize int

const (
	ForUpdate      = Lock("FOR UPDATE")
	ForNoKeyUpdate = Lock("FOR NO KEY UPDATE")
	ForShare       = Lock("FOR SHARE")
	ForKeyShare    = Lock("FOR KEY SHARE")
)

func (l Lock) String() string {
	return string(l)
}

// SELECT [ ALL | DISTINCT [ ON ( expression [, ...] ) ] ]
//    [ * | expression [ [ AS ] output_name ] [, ...] ]
//    [ FROM from_item [, ...] ]
//    [ WHERE condition ]
//    [ GROUP BY grouping_element [, ...] ]
//    [ HAVING condition [, ...] ]
//    [ WINDOW window_name AS ( window_definition ) [, ...] ]
//    [ { UNION | INTERSECT | EXCEPT } [ ALL | DISTINCT ] select ]
//    [ ORDER BY expression [ ASC | DESC | USING operator ] [ NULLS { FIRST | LAST } ] [, ...] ]
//    [ LIMIT { count | ALL } ]
//    [ OFFSET start [ ROW | ROWS ] ]
//    [ FETCH { FIRST | NEXT } [ count ] { ROW | ROWS } ONLY ]
//    [ FOR { UPDATE | NO KEY UPDATE | SHARE | KEY SHARE } [ OF table_name [, ...] ] [ NOWAIT | SKIP LOCKED ] [...] ]
type SelectQuery struct {
	expressions  []Expr
	fromTables   []Table
	fromSubquery Expr
	where        []Expr
	orderBy      []Expr
	lock         Lock
	joins        []string
	limit        int
	offset       int
	groupBy      Column
}

func Select(s ...Expr) SelectQuery {
	return SelectQuery{expressions: s}
}

func (q SelectQuery) FromSubquery(subquery Query, alias string) SelectQuery {
	q.fromSubquery = rawExpr("FROM (" + subquery.String() + ") AS " + alias)
	if len(q.expressions) == 0 {
		q.expressions = append(q.expressions, Wildcard)
	}
	return q
}

func (q SelectQuery) AlsoSelect(e ...Expr) SelectQuery {
	q.expressions = append(q.expressions, e...)
	return q
}

func (q SelectQuery) From(tables ...Table) SelectQuery {
	if len(tables) == 0 {
		log.Panicln("lorm: select query: empty tables list")
	}
	q.fromTables = tables
	if len(q.expressions) == 0 {
		for _, t := range q.fromTables {
			for _, c := range t.GetAllFields() {
				q.expressions = append(q.expressions, c)
			}
		}
	}
	return q
}

func (q SelectQuery) LeftJoin(t Table, cond Expr) SelectQuery {
	q.joins = append(q.joins, "LEFT JOIN "+t.String()+" ON "+cond.String())
	return q
}

func (q SelectQuery) InnerJoin(t Table, cond Expr) SelectQuery {
	q.joins = append(q.joins, "INNER JOIN "+t.String()+" ON "+cond.String())
	return q
}

func (q SelectQuery) Where(v ...Expr) SelectQuery   { q.where = nonEmptyExpr(v); return q }
func (q SelectQuery) OrderBy(v ...Expr) SelectQuery { q.orderBy = nonEmptyExpr(v); return q }

func (q SelectQuery) GroupBy(c Column) SelectQuery { q.groupBy = c; return q }
func (q SelectQuery) Limit(v int) SelectQuery      { q.limit = v; return q }
func (q SelectQuery) Offset(v int) SelectQuery     { q.offset = v; return q }

func (q SelectQuery) Last(c Column) SelectQuery   { return q.OrderBy(c.Desc()).Limit(1) }
func (q SelectQuery) Lock(v Lock) SelectQuery     { q.lock = v; return q }
func (q SelectQuery) ForUpdate() SelectQuery      { q.lock = ForUpdate; return q }
func (q SelectQuery) ForNoKeyUpdate() SelectQuery { q.lock = ForNoKeyUpdate; return q }

func (q SelectQuery) String() string { return q.Query() }

func (q SelectQuery) As(alias Column) Expr {
	return rawExpr("(" + q.String() + ") AS " + alias.String())
}

func (q SelectQuery) Query() string {
	if len(q.expressions) == 0 && len(q.fromTables) == 0 && q.fromSubquery == nil {
		log.Panic("invalid query expressions: empty expression, empty tables list, empty fromSubquery:", compactJSON(q))
	}

	var b strings.Builder
	b.Grow(selectQueryMaxSize)

	b.WriteString("SELECT ")
	joinExpr(&b, q.expressions, ", ")

	if q.fromSubquery != nil {
		b.WriteString(" ")
		b.WriteString(q.fromSubquery.String())
	} else if len(q.fromTables) > 0 {
		b.WriteString(" FROM ")
		for i, t := range q.fromTables {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(t.String())
		}
	}

	for _, j := range q.joins {
		b.WriteString(" ")
		b.WriteString(j)
	}

	if len(q.where) > 0 {
		b.WriteString(" WHERE (")
		for i, cond := range q.where {
			if i > 0 {
				b.WriteString(" AND ")
			}
			b.WriteString(cond.String())
		}
		b.WriteString(")")
	}

	if v := q.groupBy.String(); v != "" {
		b.WriteString(" GROUP BY ")
		b.WriteString(v)
	}

	if len(q.orderBy) > 0 {
		b.WriteString(" ORDER BY ")
		for i, v := range q.orderBy {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(v.String())
		}
	}

	if q.limit > 0 {
		b.WriteString(" LIMIT ")
		b.WriteString(strconv.Itoa(q.limit))
	}

	if q.offset > 0 {
		b.WriteString(" OFFSET ")
		b.WriteString(strconv.Itoa(q.offset))
	}

	if q.lock != "" {
		b.WriteString(" ")
		b.WriteString(q.lock.String())
	}

	return maybeGrow(b.String(), &selectQueryMaxSize)
}

func compactJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		log.Panicln(errors.Wrap(err, "json marshall fail"))
	}
	return string(data)
}
