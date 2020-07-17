package lorm

import "sync"

type BaseCache struct {
	sync.Mutex
	Name string
}
