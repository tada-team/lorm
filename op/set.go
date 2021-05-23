package op

import (
	"sort"
	"strings"
)

type Set map[Column]Expr

type setItem struct {
	Column Column
	Expr   Expr
}

func (s Set) Add(s2 Set) Set {
	res := make(Set)
	for k, v := range s {
		res[k] = v
	}
	for k, v := range s2 {
		res[k] = v
	}
	return res
}

func (s Set) SortedItems() []setItem {
	l := make([]setItem, 0, len(s))
	for k, v := range s {
		l = append(l, setItem{
			Column: k,
			Expr:   v,
		})
	}
	sort.Slice(l, func(i, j int) bool {
		return l[i].Column < l[j].Column
	})
	return l
}

func (s Set) String() string {
	res := make([]string, 0, len(s))
	for _, item := range s.SortedItems() {
		res = append(res, item.Column.BareName().Equal(item.Expr).String())
	}
	return strings.Join(res, ", ")
}
