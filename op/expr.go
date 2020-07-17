package op

import (
	"fmt"
	"strings"
)

type Expr interface {
	fmt.Stringer
}

func joinExpr(v []Expr, sep string) string {
	if v == nil || len(v) == 0 {
		return ""
	}
	bits := make([]string, 0, len(v))
	for _, f := range v {
		if f != nil && f.String() != "" {
			bits = append(bits, f.String())
		}
	}
	return strings.Join(bits, sep)
}

type rawExpr string

func (f rawExpr) String() string { return string(f) }
