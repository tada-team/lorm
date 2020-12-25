package op

import (
	"fmt"
	"testing"
)

func TestTextSearch(t *testing.T) {
	args := NewArgs()
	testArg := args.Next("test")
	col := Column("column")
	t.Run("VectorTextSearch (old style search)", func(t *testing.T) {
		result := col.VectorTextSearchEnglish(testArg)
		expectedString := fmt.Sprintf("%s @@ plainto_tsquery('%s', %s)", col, "english", testArg)
		if result.String() != expectedString {
			t.Errorf("%s != %s", result.String(), expectedString)
		}
	})
	t.Run("(new style search)", func(t *testing.T) {
		tsQuery := testArg.PhraseToTsQuery("russian")
		result := col.TsVector("russian").Search(tsQuery)
		expected := "to_tsvector('russian', column) @@ phraseto_tsquery('russian', $1)"
		if result.String() != expected {
			t.Errorf("%s != %s", result, expected)
		}
	})
}
