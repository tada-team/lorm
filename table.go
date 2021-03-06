package lorm

import (
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

	cachedFields     *[]op.Column
	cachedFieldsExpr *op.Expr
}

func NewBaseTable(name, aliasSeed string, fields ...string) BaseTable {
	return BaseTable{
		name:   name,
		alias:  nextAlias(aliasSeed),
		fields: fields,
	}
}

func (t BaseTable) Pk() op.Column { return t.Column(t.fields[0]) }

func (t BaseTable) String() string { return string(t.TableName()) }

func (t BaseTable) GetAlias() string { return t.alias }

func (t *BaseTable) SetAlias(s string) {
	t.alias = s
	t.cachedFields = nil
	t.cachedFieldsExpr = nil
}

func (t BaseTable) TableName() op.TableName {
	if t.alias != "" {
		return op.TableName(t.name + " AS " + t.alias)
	}
	return op.TableName(t.name)
}

func (t BaseTable) AllFieldsExpr() op.Expr {
	if t.cachedFieldsExpr == nil {
		expr := op.Raw(strings.Join(t.fields, ", "))
		t.cachedFieldsExpr = &expr
	}
	return *t.cachedFieldsExpr
}

func (t BaseTable) Column(v string) op.Column {
	if t.alias != "" {
		return op.Column(t.alias + "." + v)
	}
	return op.Column(t.name + "." + v)
}

func (t BaseTable) GetAllFields() []op.Column {
	if t.cachedFields == nil {
		fields := make([]op.Column, len(t.fields))
		for i, f := range t.fields {
			fields[i] = t.Column(f)
		}
		t.cachedFields = &fields
	}
	return *t.cachedFields
}
