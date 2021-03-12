package lorm

import (
	"sync"
	"testing"

	"github.com/tada-team/lorm/op"
)

func TestConcurrentMapWrite(t *testing.T) {
	var tables []op.Table
	for j := 0; j < 100; j++ {
		tables = append(tables, NewBaseTable("name", "", "id", "created"))
	}

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			for _, t := range tables {
				CachedSelect(t)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

