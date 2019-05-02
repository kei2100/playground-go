package tailf

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/kei2100/playground-go/util/tailf/logger"

	"github.com/kei2100/playground-go/util/tailf/file"
	"github.com/kei2100/playground-go/util/tailf/posfile"
)

type Follower interface {
	io.ReadCloser
}

type option struct {
	positionFile posfile.PositionFile
	//rotatedFilePathPattern string
	watchRotateInterval time.Duration
	detectRotateDelay   time.Duration
}

type OptionFunc func(o *option)

func (o *option) apply(opts ...OptionFunc) {
	o.watchRotateInterval = 200 * time.Millisecond
	o.detectRotateDelay = 5 * time.Second
	for _, fn := range opts {
		fn(o)
	}
}

func Open(name string, opts ...OptionFunc) (Follower, error) {
	opt := option{}
	opt.apply(opts...)

	f, err := file.Open(name)
	if err != nil {
		return nil, fmt.Errorf("tailf: failed to open %s: %+v", name, err)
	}

	if opt.positionFile == nil {
		fi, err := f.Stat()
		if err != nil {
			return nil, fmt.Errorf("tailf: failed to get FileInfo %s: %+v", name, err)
		}
		return &follower{
			file:                f,
			positionFile:        posfile.NewMemoryPositionFile(fi, 0),
			watchRotateInterval: opt.watchRotateInterval,
			detectRotateDelay:   opt.detectRotateDelay,
		}, nil
	}
	// TODO
	return nil, nil
}

type follower struct {
	file                *os.File
	positionFile        posfile.PositionFile
	watchRotateInterval time.Duration
	detectRotateDelay   time.Duration

	rotated <-chan struct{}
	done    chan struct{}
}

func (fw *follower) Read(p []byte) (n int, err error) {
	if fw.done == nil {
		fw.done = make(chan struct{})
	}
	if fw.rotated == nil {
		fw.rotated = watchRotate(fw.done, fw.file, fw.watchRotateInterval, fw.detectRotateDelay)
	}

	select {
	default:
		n, err := fw.file.Read(p)
		fw.positionFile.IncreaseOffset(n)
		return n, err

	case <-fw.rotated:
		if err := fw.file.Close(); err != nil {
			return 0, fmt.Errorf("tailf: an error occurred while closing target file: %+v", err)
		}
		f, err := file.Open(fw.file.Name())
		if err != nil {
			return 0, fmt.Errorf("tailf: failed to open %s: %+v", fw.file.Name(), err)
		}
		fi, err := f.Stat()
		if err != nil {
			return 0, fmt.Errorf("tailf: failed to get FileInfo %s: %+v", f.Name(), err)
		}

		fw.file = f
		fw.positionFile.Update(fi, 0)
		fw.rotated = watchRotate(fw.done, fw.file, fw.watchRotateInterval, fw.detectRotateDelay)
		return fw.Read(p)
	}
}

func (fw *follower) Close() error {
	if fw.done != nil {
		close(fw.done)
	}
	if err := fw.file.Close(); err != nil {
		return fmt.Errorf("tailf: an error occurred while closing target file: %+v", err)
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
					logger.Printf("tailf: failed to get FileInfo %s on watchRotate: %+v", file.Name(), err)
					continue
				}
				currentInfo, err := os.Stat(file.Name())
				if err != nil {
					if os.IsNotExist(err) {
						continue
					}
					logger.Printf("tailf: failed to get current FileInfo %s on watchRotate: %+v", file.Name(), err)
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
