package op

import "fmt"

type Query interface {
	fmt.Stringer
	Query() string
}
