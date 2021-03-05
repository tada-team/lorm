package op

import (
	"log"
	"testing"
)

// BenchmarkQuery/fmt.Sprintf-2         	 1326040	      1148 ns/op	     264 B/op	      13 allocs/op
// BenchmarkQuery/strings.Builder-2    	 2431648	       488 ns/op	     168 B/op	       7 allocs/op
//
// grow:
// BenchmarkQuery/strings.Builder-12         	 3883315	       281.5 ns/op	     192 B/op	       5 allocs/op
//
// bJoin:
// BenchmarkQuery/strings.Builder-12         	 3904591	       279.4 ns/op	     224 B/op	       4 allocs/op
//
// drop: joinTableNames
// BenchmarkQuery/strings.Builder-12         	 4704272	       242.9 ns/op	     208 B/op	       3 allocs/op
func BenchmarkQuery(b *testing.B) {

	b.Run("strings.Builder", func(b *testing.B) {
		b.ReportAllocs()

		var q Query
		b.Run("make query", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				q = Select(
					Wildcard,
				).From(
					tableAlias{name: "xxx"},
				).Where(
					Or(
						Raw("id = 42"),
						Raw("field = 'ttt"),
					),
					Or(
						Raw("xxxx = 42"),
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
				)
			}
		})

		b.Run("query", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				q.Query()
			}
		})
	})
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
