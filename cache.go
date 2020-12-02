package lorm

import "sync"

type BaseCache struct {
	sync.RWMutex
	Name string
}
