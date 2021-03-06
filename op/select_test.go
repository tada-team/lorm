package op

import (
	"log"
	"testing"
)

// BenchmarkQuery-12    	  416414	      2814 ns/op	    1880 B/op	      37 allocs/op
// remove fmt.Sprintf()
// BenchmarkQuery-12    	  621747	      1993 ns/op	    1480 B/op	      32 allocs/op
func BenchmarkQuery(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := With("t1", Select(One))
		w = w.With("t2", Select(Wildcard).From(tableAlias{name: "yyy"}))
		q := w.Do(Select(
			Wildcard,
		).From(
			tableAlias{name: "xxx"},
		).Where(
			Or(
				Raw("id = 42"),
				Raw("field = 'ttt"),
			),
			Or(
				Raw("field2 = 42"),
			),
		).OrderBy(
			Column("id"),
			Column("field"),
		).LeftJoin(
			tableAlias{name: "xxx"},
			And(
				Raw("id = 42"),
				Raw("field = 'ttt"),
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

func TestQueryTest(t *testing.T) {
	t.Run("test SelectQuery constructor", func(t *testing.T) {
		query := Select().Where(
			Raw("A"),
			Raw("B"),
			Raw("C"),
			Raw("D"),
		)
		answer := Raw("(A AND B AND C AND D)")
		if query.whereCondition != answer {
			t.Error("Wrong. want:", answer, "got:", query.whereCondition)
		}
	})

	t.Run("test RawQuery constructor", func(t *testing.T) {
		raw := RawQuery(
			Raw("A"),
			Raw("B"),
			Raw("C"),
			Raw("D"),
		)
		log.Println("RESULT ", raw)
		answer := rawQuery("A B C D")
		if raw != answer {
			t.Error("Wrong. want:", answer, "got:", raw)
		}
	})
}
