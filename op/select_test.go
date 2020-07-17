package op

import (
	"log"
	"testing"
)

//// BenchmarkQuery/fmt.Sprinf-2         	 1326040	      1148 ns/op	     264 B/op	      13 allocs/op
//// BenchmarkQuery/strings.Builder-2    	 2431648	       488 ns/op	     168 B/op	       7 allocs/op
//func BenchmarkQuery(b *testing.B) {
//	q := Select(Wildcard).From(tableAlias{name: "xxx"}).Where(Raw("id = 42")).OrderBy(Raw("id"))
//
//	if q.Query() != q.QueryAlt() {
//		b.Fatalf("`%s` != `%s`", q.Query(), q.QueryAlt())
//	}
//
//	b.Run("fmt.Sprinf", func(b *testing.B) {
//		b.ReportAllocs()
//		for i := 0; i < b.N; i++ {
//			q.Query()
//		}
//	})
//
//	b.Run("strings.Builder", func(b *testing.B) {
//		b.ReportAllocs()
//		for i := 0; i < b.N; i++ {
//			q.QueryAlt()
//		}
//	})
//}

func TestQueryTest(t *testing.T) {
	t.Run("test SelectQuery construtor", func(t *testing.T) {
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

	t.Run("test RawQuery construtor", func(t *testing.T) {
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
