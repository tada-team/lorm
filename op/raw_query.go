package op

import (
	"fmt"
	"strings"
)

type rawQuery string

func RawQuery(v ...interface{}) rawQuery {
	return rawQuery(strings.TrimSpace(fmt.Sprintln(v...)))
}

func (q rawQuery) String() string { return q.Query() }
func (q rawQuery) Query() string  { return string(q) }
