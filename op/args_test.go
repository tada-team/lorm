package op

import (
	"testing"
)

func TestArgs(t *testing.T) {
	args := NewArgs()
	query := Select(args.Next(42)).Query()
	if len(args) != 1 || args[0] != 42 {
		t.Error("want 1 arg, got")
	}
	if query != "SELECT $1" {
		t.Error("invalid query:", query)
	}
}
