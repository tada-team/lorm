package op

import (
	"testing"
)

// BenchmarkQuery-12    	  448762	      2611 ns/op	    1704 B/op	      50 allocs/op
func BenchmarkQuery(b *testing.B) {
	b.ReportAllocs()

	idCol := Column("id")
	ffCol := Column("ff")
	for i := 0; i < b.N; i++ {
		args := NewArgs()

		w := With("t1", Select(One))
		w = w.With("t2", Select(Wildcard).From(tableAlias{name: "yyy"}))

		q := w.Do(Select(
			Wildcard,
		).From(
			tableAlias{name: "xxx"},
		).Where(
			Or(
				idCol.Gt(args.Next(42)),
				ffCol.Equal(args.Next("ttt")),
			),
			Or(
				Raw("field2 = 42"),
			),
		).OrderBy(
			idCol,
			ffCol,
		).LeftJoin(
			tableAlias{name: "xxx"},
			And(
				idCol.Gt(args.Next(43)),
				ffCol.Equal(args.Next("ttt123123")),
			),
		).Limit(
			1,
		).Offset(
			2,
		))

		if v := q.Query(); v == "" {
			b.Fatal("empty result")
		}
	}
}

func TestSelect(t *testing.T) {
	t.Run("test SelectQuery constructor", func(t *testing.T) {
		query := Select(Wildcard).From(
			tableAlias{name: "xxx"},
		).Where(
			Raw("A"),
			Raw("B"),
			Raw("C"),
			Raw("D"),
		)
		answer := Raw("SELECT * FROM xxx WHERE (A AND B AND C AND D)")
		if query.String() != answer.String() {
			t.Error("Wrong. want:", answer, "got:", query.String())
		}
	})

	t.Run("test RawQuery constructor", func(t *testing.T) {
		raw := RawQuery(
			Raw("A"),
			Raw("B"),
			Raw("C"),
			Raw("D"),
		)
		answer := rawQuery("A B C D")
		if raw != answer {
			t.Error("Wrong. want:", answer, "got:", raw)
		}
	})
}
