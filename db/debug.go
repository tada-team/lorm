package db

import (
	"fmt"
	"path"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/tada-team/lorm/op"
)

var Debug struct {
	ShowSql bool
}

type beforeQueryHandler func(tx *Tx, qNum int64, q string, v op.Args)
type afterQueryHandler func(tx *Tx, qNum int64, q string, v op.Args, dur time.Duration)

var beforeQueryHandlers = make([]beforeQueryHandler, 0)
var afterQueryHandlers = make([]afterQueryHandler, 0)

var qNum int64

func QueryCounter() int64               { return qNum }
func BeforeQuery(fn beforeQueryHandler) { beforeQueryHandlers = append(beforeQueryHandlers, fn) }
func AfterQuery(fn afterQueryHandler)   { afterQueryHandlers = append(afterQueryHandlers, fn) }

func trackQuery(tx *Tx, q string, v op.Args) func() {
	atomic.AddInt64(&qNum, 1)
	num := qNum
	for _, fn := range beforeQueryHandlers {
		fn(tx, num, q, v)
	}
	start := time.Now()
	return func() {
		dur := time.Since(start)
		for _, fn := range afterQueryHandlers {
			fn(tx, num, q, v, dur)
		}
	}
}

func breadcrumb() string {
	_, file, no, _ := runtime.Caller(2)
	return fmt.Sprintf("%s:%d", path.Base(file), no)
}
