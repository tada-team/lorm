package lorm

import "fmt"

func isEmpty(v interface{}) bool {
	switch v.(type) {
	case int:
		return v.(int) == 0
	case int64:
		return v.(int64) == 0
	case string:
		return v.(string) == ""
	case fmt.Stringer:
		return v.(fmt.Stringer).String() == ""
	default:
		switch fmt.Sprintf("%v", v) {
		case "", "0":
			return true
		}
	}
	return false
}
