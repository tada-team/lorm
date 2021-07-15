package op

import "sync/atomic"

const maxGrowth = 5 * 1024

func maybeGrow(s string, mx *int32) string {
	n := int32(len(s))
	if mxValue := atomic.LoadInt32(mx); mxValue < maxGrowth && n > *mx {
		atomic.StoreInt32(mx, n)
	}
	return s
}
