package op

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/lib/pq"
)

var argsMaxSize int32

type Args []interface{}

func NewArgs() Args { return make(Args, 0, 8) }

func (args Args) String() string {
	var b strings.Builder
	b.Grow(int(atomic.LoadInt32(&argsMaxSize)))
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
	return Placeholder("$" + strconv.Itoa(len(*args)))
}

type ArrayMask Placeholder

func (args *Args) NextArray(v interface{}) ArrayMask {
	*args = append(*args, pq.Array(v))
	return ArrayMask("$" + strconv.Itoa(len(*args)))
}

func (args *Args) Clone() Args {
	return append(NewArgs(), *args...)
}
