package op

import (
	"fmt"
	"strings"
)

var escapeSpecialReplacer = strings.NewReplacer(
	"%", "\\%",
	"_", "\\_",
	".", "\\.",
	"*", "\\*",
)

func EscapeSpecialSymbols(s string) string {
	return escapeSpecialReplacer.Replace(s)
}

func Add(f Expr, v Expr) Expr               { return rawExpr(f.String() + " + " + v.String()) }
func Mul(f Expr, v Expr) Expr               { return rawExpr(f.String() + " * " + v.String()) }
func Aliased(alias string, f Column) Column { return Column(alias + "." + f.BareName().String()) }
func And(v ...Expr) Expr                    { return braces(v, " AND ") }
func Any(f Expr, v ArrayMask) Expr          { return rawExpr(f.String() + " = ANY(" + string(v) + ")") }
func Asc(f Expr) Expr                       { return rawExpr(f.String() + " ASC") }
func Coalesce(a, b Expr) Expr               { return rawExpr("COALESCE(" + a.String() + ", " + b.String() + ")") }
func Count(v Expr) Expr                     { return rawExpr("COUNT(" + v.String() + ")") }
func Desc(f Expr) Expr                      { return rawExpr(f.String() + " DESC") }
func Distinct(f Expr) Expr                  { return rawExpr("DISTINCT " + f.String()) }
func Exists(q Query) Expr                   { return rawExpr("EXISTS " + Subquery(q).String()) }
func FirstValue(v Expr) Expr                { return rawExpr("first_value(" + v.String() + ")") }
func GreaterThan(f Expr, v Expr) Expr       { return rawExpr(f.String() + " > " + v.String()) }
func Greatest(f Expr, v Expr) Expr          { return rawExpr("GREATEST(" + f.String() + ", " + v.String() + ")") }
func IsNotNull(f Expr) Expr                 { return rawExpr(f.String() + " IS NOT NULL") }
func IsNull(f Expr) Expr                    { return rawExpr(f.String() + " IS NULL") }
func Lag(v Expr) Expr                       { return rawExpr("LAG(" + v.String() + ")") }
func Lead(v Expr) Expr                      { return rawExpr("LEAD(" + v.String() + ")") }
func LessThan(f Expr, v Expr) Expr          { return rawExpr(f.String() + " < " + v.String()) }
func Not(f Expr) Expr                       { return rawExpr("NOT " + f.String()) }
func NotEqual(f Expr, v Expr) Expr          { return rawExpr(f.String() + " != " + v.String()) }
func NotInSubquery(f Expr, q Query) Expr    { return rawExpr(f.String() + " NOT IN (" + q.String() + ")") }
func Or(v ...Expr) Expr                     { return braces(v, " OR ") }
func Over(v ...interface{}) Expr            { return rawExpr("OVER (" + Raw(v...).String() + ")") }
func Sub(f Expr, v Expr) Expr               { return rawExpr("(" + f.String() + " - " + v.String() + ")") }
func Subquery(q Query) Expr                 { return rawExpr("(" + q.Query() + ")") }
func Sum(v Expr) Expr                       { return rawExpr("SUM(" + v.String() + ")") }

func Raw(v ...interface{}) Expr {
	if len(v) == 1 {
		s, ok := v[0].(string)
		if ok {
			return rawExpr(s)
		}
	}
	return rawExpr(strings.TrimSpace(fmt.Sprintln(v...)))
}

func PgAdvisoryXactLock(k Expr) Expr {
	return rawExpr("pg_advisory_xact_lock(" + k.String() + ")")
}

func PgAdvisoryXactLock2(k1, k2 int) Expr {
	return rawExpr(fmt.Sprintf("pg_advisory_xact_lock(%d, %d)", k1, k2))
}

func equal(f Expr, v Expr) Expr              { return rawExpr(f.String() + " = " + v.String()) }
func greaterThanOrEqual(f Expr, v Expr) Expr { return rawExpr(f.String() + " >= " + v.String()) }
func iLike(f Expr, v Expr) Expr              { return rawExpr(f.String() + " ILIKE " + v.String()) }
func inSubquery(f Expr, q Query) Expr        { return rawExpr(f.String() + " IN " + Subquery(q).String()) }
func lessThanOrEqual(f Expr, v Expr) Expr    { return rawExpr(f.String() + " <= " + v.String()) }
func notAny(f Expr, v ArrayMask) Expr {
	return rawExpr("NOT " + f.String() + " = ANY(" + string(v) + ")")
}

var bracesMaxSize int

func braces(v []Expr, sep string) Expr {
	switch len(v) {
	case 0:
		return rawExpr("")
	case 1:
		return v[0]
	default:
		var b strings.Builder
		b.Grow(bracesMaxSize)
		b.WriteString("(")
		joinExpr(&b, v, sep)
		b.WriteString(")")
		return rawExpr(maybeGrow(b.String(), &bracesMaxSize))
	}
}

func ToTsVector(lang string, f Expr) TsVector {
	return TsVector(fmt.Sprintf("to_tsvector('%s', %s)", lang, f))
}

func PlainToTsQuery(lang string, arg Expr) TsQuery {
	return TsQuery(fmt.Sprintf("plainto_tsquery('%s', %s)", lang, arg))
}

func PhraseToTsQuery(lang string, arg Expr) TsQuery {
	return TsQuery(fmt.Sprintf("phraseto_tsquery('%s', %s)", lang, arg))
}

func TextSearch(lang string, f Expr, v Expr) Expr {
	return Raw(ToTsVector(lang, f), "@@", PlainToTsQuery(lang, v))
}

func VectorTextSearch(lang string, f Expr, v Expr) Expr {
	return Raw(f, "@@", PlainToTsQuery(lang, v))
}

func Case(cond, t, f Expr) Expr {
	return rawExpr("(CASE WHEN " + cond.String() + " THEN " + t.String() +  " ELSE " + f.String() +  " END)")
}

func Union(v ...Expr) rawQuery {
	var b strings.Builder
	for i, c := range v {
		if i > 0 {
			b.WriteString(" UNION ")
		}
		b.WriteString(c.String())
	}
	return rawQuery(b.String())
}

func List(v ...Expr) Expr {
	switch len(v) {
	case 0:
		return rawExpr("")
	case 1:
		return v[0]
	default:
		var b strings.Builder
		joinExpr(&b, v, ", ")
		return rawExpr(b.String())
	}
}

func HasPrefix(v string) string { return EscapeSpecialSymbols(v) + "%" }
func HasSuffix(v string) string { return "%" + EscapeSpecialSymbols(v) }
func Contains(v string) string  { return "%" + EscapeSpecialSymbols(v) + "%" }
