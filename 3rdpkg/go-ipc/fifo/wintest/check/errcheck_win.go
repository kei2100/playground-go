//+build windows

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
	rfifo, err := fifo.New("rfifo", os.O_CREATE|os.O_RDONLY|os.O_TRUNC|fifo.O_NONBLOCK, 0666)
	if err != nil {
			log.Printf("rfifo errno error: %T", err)
			log.Printf("rfifo errno error: %v", err)
	} else {
		log.Println("rfifo close")
		rfifo.Close()
	}

	wfifo, err := fifo.New("wfifo", os.O_CREATE|os.O_WRONLY|os.O_TRUNC|fifo.O_NONBLOCK, 0666)
	if err != nil {
		// linux
		if err, ok := errors.Cause(err).(*os.PathError); ok {
			log.Printf("PathError.Err: %T", err.Err)
			log.Printf("PathError.Err: %v", err.Err)
			log.Printf("PathError.Err: %v", err.Err == syscall.ENXIO)
		}
		// win
		if err, ok := errors.Cause(err).(syscall.Errno); ok {
			log.Printf("errno error: %T", err)
			log.Printf("errno error: %v", err)
			log.Printf("errno error: %v", err == syscall.ERROR_FILE_NOT_FOUND)
		}
		log.Printf("err: %T", err)
		log.Fatalf("err: %v", err)
		log.Printf("err.Cause: %T", errors.Cause(err))
	}
	wfifo.Close()
}
