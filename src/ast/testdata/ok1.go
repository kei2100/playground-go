package testdata

import "errors"

func ok1(err error) bool {
	type fooErr struct {
		error
	}
	var fe *fooErr
	return errors.As(err, &fe)
}
