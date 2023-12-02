package testdata

import (
	"context"

	"github.com/kei2100/playground-go/src/ast/testdata/errors"
)

func ok2(err error) bool {
	return errors.Is(err, context.Canceled)
}
