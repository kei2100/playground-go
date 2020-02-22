package errors

import (
	"strings"
	"sync"
)

// errList is goroutine safe error list
type errList struct {
	errs []error
	mu   sync.RWMutex
}

// Append an err to the errList
func (el *errList) Append(err error) {
	el.mu.Lock()
	defer el.mu.Unlock()

	el.errs = append(el.errs, err)
}

// Len returns length of the errList
func (el *errList) Len() int {
	el.mu.RLock()
	defer el.mu.RUnlock()

	return len(el.errs)
}

// Join messages of the errList
func (el *errList) Join(separator string) string {
	el.mu.RLock()
	defer el.mu.RUnlock()

	var buf strings.Builder
	for i, err := range el.errs {
		if i > 0 {
			buf.WriteString(separator)
		}
		buf.WriteString(err.Error())
	}
	return buf.String()
}
