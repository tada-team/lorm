package lorm

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/tada-team/lorm/op"
)

var aliases = make(map[string]int)

func nextAlias(className string) string {
	short := ""
	for _, r := range className {
		if unicode.IsUpper(r) && unicode.IsLetter(r) {
			short += string(r)
		}
	}
	short = strings.ToLower(short)
	if short == "" {
		short = "t"
	}
	aliases[short]++
	if aliases[short] > 1 {
		short += strconv.Itoa(aliases[short])
	}
	return short
}

type BaseTable struct {
	name   string
	alias  string
	fields []string
}

func NewBaseTable(name, alias string, fields ...string) BaseTable {
	return BaseTable{
		name:   name,
		alias:  nextAlias(alias),
		fields: fields,
	}
}

func (t BaseTable) Pk() op.Column      { return t.Column(t.fields[0]) }
func (t BaseTable) GetAlias() string   { return t.alias }
func (t *BaseTable) SetAlias(s string) { t.alias = s }

func (t BaseTable) TableName() op.TableName {
	if t.alias != "" {
		return op.TableName(fmt.Sprintf("%s AS %s", t.name, t.alias))
	}
	return op.TableName(t.name)
}

func (t BaseTable) AllFieldsExpr() op.Expr {
	bits := make([]string, len(t.fields))
	for i, f := range t.GetAllFields() {
		bits[i] = f.String()
	}
	return op.Raw(strings.Join(bits, ", "))
}

func (t BaseTable) Column(v string) op.Column {
	if t.alias != "" {
		return op.Column(t.alias + "." + v)
	}
	return op.Column(t.name + "." + v)
}

func (t BaseTable) GetAllFields() []op.Column {
	fields := make([]op.Column, len(t.fields))
	for i := range t.fields {
		fields[i] = t.Column(t.fields[i])
	}
	return fields
}

func (t BaseTable) String() string {
	return string(t.TableName())
}
