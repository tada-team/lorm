package op

import (
	"fmt"
	"testing"
)

func TestTextSearchVector(t *testing.T) {
	args := NewArgs()
	testArg := args.Next("test")
	col := Column("column")
	result := col.VectorTextSearchEnglish(testArg)
	expectedString := fmt.Sprintf("%s @@ plainto_tsquery('%s', %s)", col, "english", testArg)
	if result.String() != expectedString {
		t.Errorf("%s != %s", result.String(), expectedString)
	}
}
