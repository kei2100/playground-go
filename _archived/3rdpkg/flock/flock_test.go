package flock

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gofrs/flock"
)

func TestFlock(t *testing.T) {
	t.SkipNow()

	log.Printf("start %d", os.Getpid())
	defer log.Printf("stop %d", os.Getpid())

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	fl := flock.New(filepath.Join(os.TempDir(), "go-flock.lock"))
	ok, err := fl.TryLockContext(ctx, 100*time.Millisecond)

	if !ok {
		if err == ctx.Err() {
			log.Printf("wait lock canceled %d", os.Getpid())
			return
		}
		log.Printf("unexpected error occurred while waiting for lock %d: %+v", os.Getpid(), err)
		return
	}

	defer fl.Unlock()
	time.Sleep(3 * time.Second)
}
