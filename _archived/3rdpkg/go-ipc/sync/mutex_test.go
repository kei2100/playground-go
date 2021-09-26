package sync

import (
	"log"
	"os"
	"testing"
	"time"

	ipcsync "bitbucket.org/avd/go-ipc/sync"
)

func TestMutex1(t *testing.T) {
	//t.SkipNow()

	// Lock取得したプロセスがUnlockしないで終了すると次からLock取得できなくなる。DestroyMutexすれば取得できるようになる。
	// 他のプロセスがLockを取得していても、DestroyMutexはできてしまう。さらに同名で再ロックすることができてしまう。
	if err := ipcsync.DestroyMutex("1"); err != nil {
		t.Fatal(err)
	}

	log.Printf("start %d", os.Getpid())

	mu, err := ipcsync.NewMutex("1", os.O_CREATE, 0600)
	if err != nil {
		t.Fatal(err)
	}
	ok := mu.LockTimeout(1 * time.Second)
	if ok {
		defer mu.Unlock()
		log.Printf("acquire lock %d", os.Getpid())
	} else {
		log.Printf("timeout acquire lock %d", os.Getpid())
	}

	time.Sleep(5 * time.Second)
	log.Printf("stop %d", os.Getpid())
}
