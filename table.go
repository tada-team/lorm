package lorm

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/tada-team/lorm/op"
)

var (
	aliases      = make(map[string]int)
	defaultAlias = "t"
)

func nextAlias(className string) string {
	var b strings.Builder
	for _, r := range className {
		if unicode.IsUpper(r) && unicode.IsLetter(r) {
			b.WriteRune(r)
		}
	}

	short := strings.ToLower(b.String())
	if short == "" {
		short = defaultAlias
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

	cachedPk         op.Column
	cachedTableName  string
	cachedFields     *[]op.Column
	cachedFieldsExpr *op.Expr
	cachedColumns    map[string]op.Column
}

func NewBaseTable(name, aliasSeed string, fields ...string) BaseTable {
	t := BaseTable{name: name, fields: fields}
	t.SetAlias(nextAlias(aliasSeed))
	return t
}

func (t BaseTable) Pk() op.Column             { return t.cachedPk }
func (t BaseTable) String() string            { return t.cachedTableName }
func (t BaseTable) GetAlias() string          { return t.alias }
func (t BaseTable) AllFieldsExpr() op.Expr    { return *t.cachedFieldsExpr }
func (t BaseTable) GetAllFields() []op.Column { return *t.cachedFields }

func (t *BaseTable) SetAlias(s string) {
	n := len(t.fields)

	t.alias = s
	t.cachedColumns = make(map[string]op.Column, n)

	var b strings.Builder
	size := 0
	for _, f := range t.fields {
		size += len(f)
	}
	if t.alias != "" {
		size = (len(t.alias) + 1 + 2) * n
	} else {
		size = (len(t.name) + 1 + 2) * n
	}
	b.Grow(size)

	fields := make([]op.Column, n)
	for i, f := range t.fields {
		c := t.Column(f)
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(c.String())
		fields[i] = c
	}

	expr := op.Raw(b.String())
	t.cachedFieldsExpr = &expr
	t.cachedFields = &fields
	t.cachedPk = t.Column(t.fields[0])

	t.cachedTableName = t.name
	if t.alias != "" {
		t.cachedTableName += " AS " + t.alias
	}
}

func (t BaseTable) Column(v string) op.Column {
	res := t.cachedColumns[v]
	if res == "" {
		if t.alias != "" {
			res = op.Column(t.alias + "." + v)
		} else {
			res = op.Column(t.name + "." + v)
		}
		t.cachedColumns[v] = res
	}
	return res
}
