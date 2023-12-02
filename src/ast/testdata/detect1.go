package testdata

import (
	"context"
	"errors"
)

func detect1(err error) bool {
	return errors.Is(err, context.Canceled)
}
