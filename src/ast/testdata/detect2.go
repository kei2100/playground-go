package testdata

import (
	"context"
	stdliberrors "errors"
)

func detect2(err error) bool {
	return stdliberrors.Is(err, context.Canceled)
}
