package go1_20

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextWithCancelCause(t *testing.T) {
	myError := errors.New("my error")
	ctx, can := context.WithCancelCause(context.Background())
	can(myError)

	assert.Equal(t, "context canceled", ctx.Err().Error())
	assert.Equal(t, "my error", context.Cause(ctx).Error())
}
