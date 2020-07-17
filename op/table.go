package op

import (
	"fmt"
	"strings"
)

type Table interface {
	fmt.Stringer
	TableName() TableName
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
func (t tableAlias) TableName() TableName   { return t.name }
func (t tableAlias) GetAllFields() []Column { panic("not implemented") }
func (t tableAlias) Pk() Column             { panic("not implemented") }

type TableName string

//func (t TableName) TableName() TableName { return TableName(t) }
//func (t TableName) String() string    { return t.TableName() }
//func (t TableName) As(v string) Table { return Table(fmt.Sprintf("%s AS %s", t, v)) }

func joinTableNames(v []Table, sep string) string {
	bits := make([]string, 0)
	for _, t := range v {
		bits = append(bits, t.String())
	}
	return strings.Join(bits, sep)
}
