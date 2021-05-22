package op

import (
	"testing"
)

// BenchmarkInsertQuery-12    	  542629	      2166 ns/op	     936 B/op	      23 allocs/op
// ==>
// BenchmarkInsertQuery-12    	  789996	      1289 ns/op	     956 B/op	      15 allocs/op
func BenchmarkInsertQuery(b *testing.B) {
	b.ReportAllocs()

	idCol := Column("id")
	ffCol := Column("ff")
	for i := 0; i < b.N; i++ {
		args := NewArgs()
		q := Insert(tableAlias{name: "xxx"}, Set{
			idCol: args.Next(1),
			ffCol: args.Next(2),
		})
		if v := q.Query(); v == "" {
			b.Fatal("empty result")
		}
	}
}

func TestInsert(t *testing.T) {
	idCol := Column("id")
	ffCol := Column("ff")
	args := NewArgs()
	q := Insert(tableAlias{name: "xxx"}, Set{
		idCol: args.Next(1),
		ffCol: args.Next(2),
	}, Set{
		idCol: args.Next(4),
		ffCol: args.Next(5),
	})

	want := "INSERT INTO xxx (id, ff) VALUES ($1, $2), ($3, $4)"
	if v := q.Query(); v != want {
		t.Fatal("invalid result:", v, "want:", want)
	}
}
