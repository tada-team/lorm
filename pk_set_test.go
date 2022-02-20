package lorm

import "testing"

func TestSet(t *testing.T) {
	type Pk string

	set := make(Set[Pk])
	set.Add("1")

	if len(set) != 1 {
		t.Error("want size: 1, got:", len(set))
	}
}
