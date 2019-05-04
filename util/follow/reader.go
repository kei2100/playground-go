package follow

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	errAndClose := func(err error) (Reader, error) {
		if cErr := f.Close(); cErr != nil {
			logger.Printf("follow: an error occurred while closing the file %s: %+v", name, cErr)
		}
		return nil, err
	}
	fileInfo, err := f.Stat()
	if err != nil {
		return errAndClose(err)
	}

	if opt.positionFile == nil {
		positionFile := posfile.NewMemoryPositionFile(fileInfo, 0)
		return newReader(f, positionFile, opt), nil
	}
	if !os.SameFile(fileInfo, opt.positionFile.FileInfo()) {
		logger.Printf("follow: file not found that matches fileInfo of the positionFile %+v. reset positionFile.", opt.positionFile.FileInfo())
		opt.positionFile.Update(fileInfo, 0)
		return newReader(f, opt.positionFile, opt), nil
	}
	if fileInfo.Size() < opt.positionFile.Offset() {
		// consider file truncated
		logger.Printf("follow: incorrect positionFile offset %d. file size %d. reset offset to %d.", opt.positionFile.Offset(), fileInfo.Size(), fileInfo.Size())
		opt.positionFile.Update(fileInfo, fileInfo.Size())
	}
	offset, err := f.Seek(opt.positionFile.Offset(), 0)
	if err != nil {
		return errAndClose(err)
	}
	if offset != opt.positionFile.Offset() {
		return errAndClose(fmt.Errorf("follow: seems like seek failed. positionFile offset %d. file offset %d", opt.positionFile.Offset(), offset))
	}
	return newReader(f, opt.positionFile, opt), nil
}

type reader struct {
	file         *os.File
	positionFile posfile.PositionFile

	watchRotateInterval time.Duration
	detectRotateDelay   time.Duration

	closed  chan struct{}
	rotated <-chan struct{}
}

func newReader(file *os.File, positionFile posfile.PositionFile, opt option) *reader {
	closed := make(chan struct{})
	watchRotateInterval := opt.watchRotateInterval
	detectRotateDelay := opt.detectRotateDelay
	var rotated <-chan struct{}
	if opt.followRotate {
		rotated = watchRotate(closed, file, watchRotateInterval, detectRotateDelay)
	}
	return &reader{
		file:                file,
		positionFile:        positionFile,
		watchRotateInterval: watchRotateInterval,
		detectRotateDelay:   detectRotateDelay,
		closed:              closed,
		rotated:             rotated,
	}
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

func findAndOpenSameFile(fileInfo os.FileInfo, pathPattern string) (*os.File, error) {
	entries, err := filepath.Glob(pathPattern)
	if err != nil {
		return nil, err
	}
	for _, p := range entries {
		candidateInfo, err := os.Stat(p)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		if !os.SameFile(candidateInfo, fileInfo) {
			continue
		}
		sameFile, err := file.Open(p)
		if err != nil {
			if os.IsNotExist(err) {
				// sameFile renamed?
				continue
			}
			return nil, err
		}
		return sameFile, nil
	}
	return nil, os.ErrNotExist
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
