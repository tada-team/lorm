package lorm

import (
	"testing"
)

func TestBaseTable_AllFieldsExpr(t *testing.T) {
	table := NewBaseTable("name", "", "id", "created")

	if v := table.GetAlias(); v != "t" {
		t.Fatal("invalid alias:", v, "want: aaa")
	}

	if v := table.GetAllFields()[0]; v != "t.id" {
		t.Fatal("invalid field:", v, "want: t.id")
	}

	table.SetAlias("k123")

	if v := table.GetAllFields()[0]; v != "k123.id" {
		t.Fatal("invalid field:", v, "want: k123.id")
	}
}
