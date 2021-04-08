package lorm

import (
	"fmt"
	"path"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/tada-team/lorm/op"
)

type (
	beforeQueryHandler func(tx *Tx, qNum int64, q string, v op.Args)
	afterQueryHandler  func(tx *Tx, qNum int64, q string, v op.Args, dur time.Duration)
	cacheUsedHandler   func(name string)
)

var (
	qNum                int64
	beforeQueryHandlers []beforeQueryHandler
	afterQueryHandlers  []afterQueryHandler
	cacheUsedHandlers   []cacheUsedHandler
)

func QueryCounter() int64               { return qNum }
func BeforeQuery(fn beforeQueryHandler) { beforeQueryHandlers = append(beforeQueryHandlers, fn) }
func AfterQuery(fn afterQueryHandler)   { afterQueryHandlers = append(afterQueryHandlers, fn) }
func OnCacheUsed(fn cacheUsedHandler)   { cacheUsedHandlers = append(cacheUsedHandlers, fn) }

func CacheUsed(name string) {
	for _, fn := range cacheUsedHandlers {
		fn(name)
	}
}

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
