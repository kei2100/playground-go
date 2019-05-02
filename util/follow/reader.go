package follow

import (
	"io"
	"os"
	"time"

	"github.com/kei2100/playground-go/util/follow/file"
	"github.com/kei2100/playground-go/util/follow/logger"
	"github.com/kei2100/playground-go/util/follow/posfile"
)

// Reader is an interface.
type Reader interface {
	io.ReadCloser
}

// Open opens the named file for reading.
func Open(name string, opts ...OptionFunc) (Reader, error) {
	opt := option{}
	opt.apply(opts...)

	f, err := file.Open(name)
	if err != nil {
		return nil, err
	}

	if opt.positionFile == nil {
		fi, err := f.Stat()
		if err != nil {
			return nil, err
		}
		closed := make(chan struct{})
		watchRotateInterval := opt.watchRotateInterval
		detectRotateDelay := opt.detectRotateDelay
		return &reader{
			file:                f,
			positionFile:        posfile.NewMemoryPositionFile(fi, 0),
			watchRotateInterval: watchRotateInterval,
			detectRotateDelay:   detectRotateDelay,
			closed:              closed,
			rotated:             watchRotate(closed, f, watchRotateInterval, detectRotateDelay),
		}, nil
	}
	// TODO
	return nil, nil
}

type reader struct {
	file         *os.File
	positionFile posfile.PositionFile

	watchRotateInterval time.Duration
	detectRotateDelay   time.Duration

	closed  chan struct{}
	rotated <-chan struct{}
}

// Read reads up to len(b) bytes from the File.
func (r *reader) Read(p []byte) (n int, err error) {
	select {
	default:
		n, err := r.file.Read(p)
		r.positionFile.IncreaseOffset(n)
		return n, err

	case <-r.rotated:
		if err := r.file.Close(); err != nil {
			return 0, err
		}
		f, err := file.Open(r.file.Name())
		if err != nil {
			return 0, err
		}
		fi, err := f.Stat()
		if err != nil {
			return 0, err
		}
		r.file = f
		r.positionFile.Update(fi, 0)
		r.rotated = watchRotate(r.closed, r.file, r.watchRotateInterval, r.detectRotateDelay)
		return r.Read(p)
	}
}

// Close closes the follow.Reader.
func (r *reader) Close() error {
	if r.closed != nil {
		close(r.closed)
	}
	if err := r.file.Close(); err != nil {
		return err
	}
	return nil
}

func watchRotate(done chan struct{}, file *os.File, interval, notifyDelay time.Duration) (rotated <-chan struct{}) {
	notify := make(chan struct{})

	go func() {
		tick := time.NewTicker(interval)
		defer tick.Stop()
		for {
			select {
			case <-done:
				return
			case <-tick.C:
				fileInfo, err := file.Stat()
				if err != nil {
					logger.Printf("follow: failed to get FileInfo %s on watchRotate: %+v", file.Name(), err)
					continue
				}
				currentInfo, err := os.Stat(file.Name())
				if err != nil {
					if os.IsNotExist(err) {
						continue
					}
					logger.Printf("follow: failed to get current FileInfo %s on watchRotate: %+v", file.Name(), err)
					continue
				}
				if !os.SameFile(fileInfo, currentInfo) {
					<-time.After(notifyDelay)
					close(notify)
					return
				}
			}
		}
	}()

	return notify
}
