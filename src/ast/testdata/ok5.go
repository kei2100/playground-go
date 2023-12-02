package testdata

import (
	"context"

	errors "golang.org/x/xerrors"
)

func ok5(err error) bool {
	return errors.Is(err, context.Canceled)
}
