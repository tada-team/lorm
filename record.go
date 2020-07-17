package lorm

import (
	"github.com/tada-team/lorm/db"
	"github.com/tada-team/lorm/op"
)

type Record interface {
	db.Transactional
	GetAllFields() []interface{}
	HasPk() bool
	PkCond(args *op.Args) op.Expr
	Save() error
}

type Saveable interface{ Save() error }
type Deletable interface{ Delete() error }
type Reloadable interface{ Reload() error }
