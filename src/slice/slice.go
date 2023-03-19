package slice

func DeleteAt[T any](s []T, i int) []T {
	// e.g.
	// s = []int{0, 1, 2, 3, 4}
	// i = 2

	// s[i:]   = {2, 3, 4}
	// s[i+1:] = {3, 4}
	n := copy(s[i:], s[i+1:])
	// s = {0, 1, 3, 4, 4}
	// n = 2

	// s[:i+n] = {0, 1, 3, 4}
	return s[:i+n]

	// Note:
	// return append(s[:i], s[i+1:]...) でも同じ値を返せるがアロケーションが一回発生してしまう
}
