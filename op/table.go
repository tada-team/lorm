package op

import (
	"fmt"
)

type Table interface {
	fmt.Stringer
	GetAllFields() []Column
	Pk() Column
}

type tableAlias struct {
	name TableName
}

func TableAlias(name string) tableAlias {
	return tableAlias{name: TableName(name)}
}

func (t tableAlias) String() string         { return string(t.name) }
func (t tableAlias) GetAllFields() []Column { panic("not implemented") }
func (t tableAlias) Pk() Column             { panic("not implemented") }

type TableName string
