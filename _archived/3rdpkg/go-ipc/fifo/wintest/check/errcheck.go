//go:build !windows
// +build !windows

package check

import (
	"log"
	"os"
	"syscall"

	"bitbucket.org/avd/go-ipc/fifo"
	"github.com/pkg/errors"
)

// Errcheck func
func Errcheck() {
	wfifo, err := fifo.New("fifo", os.O_CREATE|os.O_WRONLY|os.O_TRUNC|fifo.O_NONBLOCK, 0666)
	if err != nil {
		// linux
		if err, ok := errors.Cause(err).(*os.PathError); ok {
			log.Printf("PathError.Err: %T", err.Err)
			log.Printf("PathError.Err: %v", err.Err)
			log.Printf("PathError.Err: %v", err.Err == syscall.ENXIO)
		}
		log.Printf("err: %T", err)
		log.Fatalf("err: %v", err)
		log.Printf("err.Cause: %T", errors.Cause(err))
	}
	wfifo.Close()
}
