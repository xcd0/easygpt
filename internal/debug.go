package internal

func debug_count_slice(ss [][]string) int {
	l := 0
	for _, s := range ss {
		l += len(s)
	}
	return l
}
