package op

import (
	"fmt"
	"strings"
)

type Column string

const (
	Wildcard = Column("*")
	Random   = Column("RANDOM()")
)

func (f Column) Add(v Expr) Expr            { return Add(f, v) }
func (f Column) Any(v ArrayMask) Expr       { return Any(f, v) }
func (f Column) As(c Column) Expr           { return rawExpr(fmt.Sprintf("(%s) AS %s", f, c)) }
func (f Column) Asc() Expr                  { return Asc(f) }
func (f Column) Desc() Expr                 { return Desc(f) }
func (f Column) Equal(v Expr) Expr          { return equal(f, v) }
func (f Column) Gt(v Expr) Expr             { return GreaterThan(f, v) }
func (f Column) Gte(v Expr) Expr            { return greaterThanOrEqual(f, v) }
func (f Column) ILike(v Expr) Expr          { return iLike(f, v) }
func (f Column) InSubquery(q Query) Expr    { return inSubquery(f, q) }
func (f Column) IsNotNull() Expr            { return IsNotNull(f) }
func (f Column) IsNull() Expr               { return IsNull(f) }
func (f Column) Lt(v Expr) Expr             { return LessThan(f, v) }
func (f Column) Lte(v Expr) Expr            { return lessThanOrEqual(f, v) }
func (f Column) NotAny(v ArrayMask) Expr    { return notAny(f, v) }
func (f Column) NotEqual(v Expr) Expr       { return NotEqual(f, v) }
func (f Column) NotInSubquery(q Query) Expr { return NotInSubquery(f, q) }
func (f Column) String() string             { return string(f) }
func (f Column) Sub(v Expr) Expr            { return Sub(f, v) }

// Use "russian_engstop" configuration
func (f Column) TextSearchRussianEngStop(v Expr) Expr { return TextSearch("russian_engstop", f, v) }
func (f Column) TextSearchEnglish(v Expr) Expr        { return TextSearch("english", f, v) }
func (f Column) TextSearchRussian(v Expr) Expr        { return TextSearch("russian", f, v) }

func (f Column) NullableNotEqual(c Column) Expr {
	return Or(
		NotEqual(f, c),
		And(f.IsNull(), c.IsNotNull()),
		And(f.IsNotNull(), c.IsNull()),
	)
}

func (f Column) BareName() Column {
	bits := strings.SplitN(f.String(), ".", 2)
	if len(bits) > 1 {
		return Column(bits[1])
	}
	return f
}
