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

func Add(f Expr, v Expr) Expr               { return rawExpr(fmt.Sprintf("(%s + %s)", f, v)) }
func Mul(f Expr, v Expr) Expr               { return rawExpr(fmt.Sprintf("(%s * %s)", f, v)) }
func Aliased(alias string, f Column) Column { return Column(alias + "." + f.BareName().String()) }
func And(v ...Expr) Expr                    { return braces(v, " AND ") }
func Any(f Expr, v ArrayMask) Expr          { return rawExpr(fmt.Sprintf("%s = ANY(%s)", f, string(v))) }
func Asc(f Expr) Expr                       { return rawExpr(fmt.Sprintf("%s ASC", f)) }
func Coalesce(a, b Expr) Expr               { return rawExpr(fmt.Sprintf("COALESCE(%s, %s)", a, b)) }
func Count(v Expr) Expr                     { return rawExpr(fmt.Sprintf("COUNT(%s)", v)) }
func Desc(f Expr) Expr                      { return rawExpr(fmt.Sprintf("%s DESC", f)) }
func Distinct(f Expr) Expr                  { return rawExpr(fmt.Sprintf("DISTINCT %s", f)) }
func Exists(q Query) Expr                   { return rawExpr(fmt.Sprintf("EXISTS %s", Subquery(q))) }
func FirstValue(v Expr) Expr                { return rawExpr(fmt.Sprintf("first_value(%s)", v)) }
func GreaterThan(f Expr, v Expr) Expr       { return rawExpr(fmt.Sprintf("%s > %s", f, v)) }
func Greatest(f Expr, v Expr) Expr          { return rawExpr(fmt.Sprintf("GREATEST(%s, %s)", f, v)) }
func IsNotNull(f Expr) Expr                 { return rawExpr(fmt.Sprintf("%s IS NOT NULL", f)) }
func IsNull(f Expr) Expr                    { return rawExpr(fmt.Sprintf("%s IS NULL", f)) }
func Lag(v Expr) Expr                       { return rawExpr(fmt.Sprintf("LAG(%s)", v)) }
func Lead(v Expr) Expr                      { return rawExpr(fmt.Sprintf("LEAD(%s)", v)) }
func LessThan(f Expr, v Expr) Expr          { return rawExpr(fmt.Sprintf("%s < %s", f, v)) }
func Not(f Expr) Expr                       { return rawExpr(fmt.Sprintf("NOT %s", f.String())) }
func NotEqual(f Expr, v Expr) Expr          { return rawExpr(fmt.Sprintf("%s != %s", f, v)) }
func NotInSubquery(f Expr, q Query) Expr    { return rawExpr(fmt.Sprintf("%s NOT IN %s", f, Subquery(q))) }
func Or(v ...Expr) Expr                     { return braces(v, " OR ") }
func Over(v ...interface{}) Expr            { return rawExpr(fmt.Sprintf("OVER (%s)", Raw(v...))) }
func Raw(v ...interface{}) Expr             { return rawExpr(strings.TrimSpace(fmt.Sprintln(v...))) }
func Sub(f Expr, v Expr) Expr               { return rawExpr(fmt.Sprintf("(%s - %s)", f, v)) }
func Subquery(q Query) Expr                 { return rawExpr("(" + q.Query() + ")") }
func Sum(v Expr) Expr                       { return rawExpr(fmt.Sprintf("SUM(%s)", v)) }

func PgAdvisoryXactLock(k Expr) Expr {
	return rawExpr(fmt.Sprintf("pg_advisory_xact_lock(%s)", k))
}

func PgAdvisoryXactLock2(k1, k2 int) Expr {
	return rawExpr(fmt.Sprintf("pg_advisory_xact_lock(%d, %d)", k1, k2))
}

func equal(f Expr, v Expr) Expr              { return rawExpr(fmt.Sprintf("%s = %s", f, v)) }
func greaterThanOrEqual(f Expr, v Expr) Expr { return rawExpr(fmt.Sprintf("%s >= %s", f, v)) }
func iLike(f Expr, v Expr) Expr              { return rawExpr(fmt.Sprintf("%s ILIKE %s", f, v)) }
func inSubquery(f Expr, q Query) Expr        { return rawExpr(fmt.Sprintf("%s IN %s", f, Subquery(q))) }
func lessThanOrEqual(f Expr, v Expr) Expr    { return rawExpr(fmt.Sprintf("%s <= %s", f, v)) }
func notAny(f Expr, v ArrayMask) Expr        { return rawExpr(fmt.Sprintf("NOT %s = ANY(%s)", f, string(v))) }

func braces(v []Expr, sep string) Expr {
	switch len(v) {
	case 0:
		return rawExpr("")
	case 1:
		return v[0]
	default:
		return rawExpr("(" + joinExpr(v, sep) + ")")
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
func VectorTextSearch(lang string, f Expr, v Expr) Expr { return Raw(f, "@@", PlainToTsQuery(lang, v)) }

func Case(cond, t, f Expr) Expr {
	return rawExpr(fmt.Sprintf("(CASE WHEN %s THEN %s ELSE %s END)", cond, t, f))
}

func Union(v ...Expr) rawQuery {
	return rawQuery(joinExpr(v, " UNION "))
}

func List(v ...Expr) Expr {
	switch len(v) {
	case 0:
		return rawExpr("")
	case 1:
		return v[0]
	default:
		return rawExpr(joinExpr(v, ", "))
	}
}

func HasPrefix(v string) string { return EscapeSpecialSymbols(v) + "%" }
func HasSuffix(v string) string { return "%" + EscapeSpecialSymbols(v) }
func Contains(v string) string  { return "%" + EscapeSpecialSymbols(v) + "%" }
