package testing

import "testing"

// Setup executes given setups, and return teardown function
func Setup(t *testing.T, setups ...func(*testing.T) (teardown func())) (teardowns func()) {
	t.Helper()

	tds := make([]func(), len(setups))
	for i, up := range setups {
		tds[len(tds)-1-i] = up(t)
	}
	return func() {
		t.Helper()
		for _, td := range tds {
			td()
		}
	}
}
