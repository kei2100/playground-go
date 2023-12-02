package testdata

import (
	"context"

	"github.com/kei2100/playground-go/src/ast/testdata/myerrors"
)

func ok3(err error) bool {
	return myerrors.Is(err, context.Canceled)
}
