package lorm

import "github.com/tada-team/lorm/op"

type Record interface {
	Transactional
	GetAllFields() []interface{}
	HasPk() bool
	PkCond(args *op.Args) op.Expr
	NewPk()
	Save() error
	//table() op.Table
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
