package testdata

import (
	"context"

	"golang.org/x/xerrors"
)

func ok4(err error) bool {
	return xerrors.Is(err, context.Canceled)
}
