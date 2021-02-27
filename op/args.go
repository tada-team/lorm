package op

import (
	"fmt"
	"strings"

	"github.com/lib/pq"
)

var argsMaxSize int

type Args []interface{}

func NewArgs() Args { return make(Args, 0) }

func (args Args) String() string {
	var b strings.Builder
	b.Grow(argsMaxSize)
	b.WriteString("sqlargs{")
	for i, v := range args {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(fmt.Sprintf(`"$%d": %v`, i+1, v))
	}
	b.WriteString("}")
	return maybeGrow(b.String(), &argsMaxSize)
}

const (
	EmptyString = rawExpr("''")
	EmptyJSON   = rawExpr("'{}'::json")
	False       = rawExpr("'false'::bool")
	Now         = rawExpr("NOW()")
	Null        = rawExpr("NULL")
	One         = rawExpr("1")
	True        = rawExpr("'true'::bool")
	Zero        = rawExpr("0")
)

func (args *Args) Next(v interface{}) Placeholder {
	*args = append(*args, v)
	return Placeholder(fmt.Sprintf("$%d", len(*args)))
}

type ArrayMask Placeholder

func (args *Args) NextArray(v interface{}) ArrayMask {
	*args = append(*args, pq.Array(v))
	return ArrayMask(fmt.Sprintf("$%d", len(*args)))
}
