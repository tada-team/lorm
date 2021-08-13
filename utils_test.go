package lorm

import "testing"

type (
	uid string
	pk  int64
)

func TestIsEmpty(t *testing.T) {
	for k, v := range map[string]interface{}{
		"int":   0,
		"int64": int64(0),
		"str":   "",
		"uid":   uid(""),
		"pk":    pk(0),
	} {
		t.Run(k, func(t *testing.T) {
			if !isEmpty(&v) {
				t.Error(k, "not empty:", v)
			}
		})
	}
	for k, v := range map[string]interface{}{
		"int":   1,
		"int64": int64(1),
		"str":   "123",
		"uid":   uid(UUID()),
		"pk":    pk(-1),
	} {
		t.Run(k, func(t *testing.T) {
			if isEmpty(&v) {
				t.Error("empty:", v)
			}
		})
	}
}
