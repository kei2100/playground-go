package errors

import "strings"

// Join concatenates the error messages of errs to create a single error message.
// The separator string sep is placed between messages in the resulting string.
func Join(errs []error, sep string) string {
	switch len(errs) {
	case 0:
		return ""
	case 1:
		return errs[0].Error()
	}
	b := strings.Builder{}
	for i := range errs {
		if i > 0 {
			b.WriteString(sep)
		}
		b.WriteString(errs[i].Error())
	}
	return b.String()
}
