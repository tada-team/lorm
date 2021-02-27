package op

import (
	"encoding/json"
	"fmt"
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
	expressions    []Expr
	fromTables     []Table
	fromSubquery   Expr
	whereCondition Expr
	orderBy        []Expr
	lock           Lock
	joins          []string
	limit          int
	offset         int
	groupBy        Column
}

func Select(s ...Expr) SelectQuery {
	return SelectQuery{expressions: s}
}

func (q SelectQuery) FromSubquery(subquery Query, alias string) SelectQuery {
	q.fromSubquery = rawExpr(fmt.Sprintf("FROM (%s) AS %s", subquery, alias))
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
		log.Panicln("select query: empty tables list")
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
	q.joins = append(q.joins, fmt.Sprintf("LEFT JOIN %s ON %s", t, cond))
	return q
}

func (q SelectQuery) InnerJoin(t Table, cond Expr) SelectQuery {
	q.joins = append(q.joins, fmt.Sprintf("INNER JOIN %s ON %s", t, cond))
	return q
}

func (q SelectQuery) Where(conds ...Expr) SelectQuery {
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

func (q SelectQuery) GroupBy(f Column) SelectQuery {
	q.groupBy = f
	return q
}

func (q SelectQuery) Limit(v int) SelectQuery       { q.limit = v; return q }
func (q SelectQuery) Offset(v int) SelectQuery      { q.offset = v; return q }
func (q SelectQuery) OrderBy(v ...Expr) SelectQuery { q.orderBy = v; return q }

func (q SelectQuery) Last(c Column) SelectQuery { return q.OrderBy(c.Desc()).Limit(1) }
func (q SelectQuery) Lock(v Lock) SelectQuery   { q.lock = v; return q }
func (q SelectQuery) ForUpdate() SelectQuery    { q.lock = ForUpdate; return q }

func (q SelectQuery) ForNoKeyUpdate() SelectQuery { q.lock = ForNoKeyUpdate; return q }

func (q SelectQuery) String() string { return q.Query() }

func (q SelectQuery) As(alias Column) Expr {
	return rawExpr(fmt.Sprintf("(%s) AS %s", q.Query(), alias))
}

func (q SelectQuery) Query() string {
	if len(q.expressions) == 0 && len(q.fromTables) == 0 && q.fromSubquery == nil {
		log.Panic("invalid query expressions: empty expression, empty tables list, empty fromSubquery:", compactJSON(q))
	}

	var b strings.Builder
	b.Grow(selectQueryMaxSize)

	b.WriteString("SELECT ")
	b.WriteString(joinExpr(q.expressions, ", "))

	if q.fromSubquery != nil {
		b.WriteString(" ")
		b.WriteString(q.fromSubquery.String())
	} else if v := joinTableNames(q.fromTables, ", "); v != "" {
		b.WriteString(" FROM ")
		b.WriteString(v)
	}

	if q.joins != nil {
		for _, j := range q.joins {
			b.WriteString(" ")
			b.WriteString(j)
		}
	}

	if q.whereCondition != nil && q.whereCondition.String() != "" {
		b.WriteString(" WHERE ")
		b.WriteString(q.whereCondition.String())
	}

	if q.groupBy.String() != "" {
		b.WriteString(" GROUP BY ")
		b.WriteString(q.groupBy.String())
	}

	//if q.window != nil && q.window.String() != "" {
	//	b.WriteString(" ")
	//	b.WriteString(q.window.String())
	//}

	if v := joinExpr(q.orderBy, ",  "); v != "" {
		b.WriteString(" ORDER BY ")
		b.WriteString(v)
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
