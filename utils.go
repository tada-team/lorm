package lorm

import (
	"encoding/json"
)

func isEmpty(v interface{}) bool {
	t, err := json.Marshal(v) // FIXME: reflect
	if err != nil {
		return false
	}
	s := string(t)
	return s == `""` || s == "0"
}
