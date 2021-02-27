package op

const maxGrowth = 5 * 1024

func maybeGrow(s string, mx *int) string {
	n := len(s)
	if *mx < maxGrowth && n > *mx {
		*mx = n
	}
	return s
}
