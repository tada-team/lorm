package lorm

import (
	"fmt"
	"sync"

	"github.com/tada-team/lorm/op"
)

type BaseRecord struct {
	BaseTransactional
	sync.RWMutex
}

type Record interface {
	Transactional
	sync.Locker
	fmt.Stringer
	GetAllFields() []interface{}
	HasPk() bool
	PkCond(args *op.Args) op.Expr
	NewPk()
	Save() error
	PreSave() error
	PostSave() error
}

type Deletable interface {
	Delete() error
}

type Saveable interface {
	Save() error
}

type Reloadable interface {
	Reload() error
}
