package testdata

import (
	"context"
	. "errors"
)

func detect3(err error) bool {
	return Is(err, context.Canceled)
}
