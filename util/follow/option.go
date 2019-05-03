package follow

import (
	"time"

	"github.com/kei2100/playground-go/util/follow/posfile"
)

type option struct {
	positionFile posfile.PositionFile
	//rotatedFilePathPattern string
	followRotate        bool
	detectRotateDelay   time.Duration
	watchRotateInterval time.Duration
}

// OptionFunc let you change follow.Reader behavior.
type OptionFunc func(o *option)

func (o *option) apply(opts ...OptionFunc) {
	o.followRotate = true
	o.watchRotateInterval = 200 * time.Millisecond
	o.detectRotateDelay = 5 * time.Second
	for _, fn := range opts {
		fn(o)
	}
}

// WithFollowRotate let you change followRotate
func WithFollowRotate(follow bool) OptionFunc {
	return func(o *option) {
		o.followRotate = follow
	}
}

// WithDetectRotateDelay let you change detectRotateDelay
func WithDetectRotateDelay(v time.Duration) OptionFunc {
	return func(o *option) {
		o.detectRotateDelay = v
	}
}

// WithWatchRotateInterval let you change watchRotateInterval
func WithWatchRotateInterval(v time.Duration) OptionFunc {
	return func(o *option) {
		o.watchRotateInterval = v
	}
}
