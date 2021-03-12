package lorm

import (
	"testing"
)

func TestBaseTable_AllFieldsExpr(t *testing.T) {
	table := NewBaseTable("name", "SomeTable", "id", "created")

	if v := table.GetAlias(); v != "st" {
		t.Fatal("invalid alias:", v, "want: st")
	}

	if v := table.GetAllFields()[0]; v != "st.id" {
		t.Fatal("invalid field:", v, "want: st.id")
	}

	table.SetAlias("k123")

	if v := table.GetAllFields()[0]; v != "k123.id" {
		t.Fatal("invalid field:", v, "want: k123.id")
	}
}
