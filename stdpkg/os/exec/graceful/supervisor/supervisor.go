package supervisor

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/kei2100/playground-go/stdpkg/os/exec/graceful/worker"
)

// Supervisor manages worker process(es)
type Supervisor struct {
	Command       string
	Args          []string
	Listeners     []net.Listener
	WaitReadyFunc func(ctx context.Context, extraFileConns []net.Conn) error

	AutoRestartEnabled bool
	AutoRestartTimeout time.Duration

	StopOldDelay time.Duration

	worker   *worker.Worker
	workerMu sync.RWMutex

	chanCloseMonitor chanCloseMonitor
}

// Start Supervisor
// blocks until the worker process is done
func (s *Supervisor) Start(ctx context.Context) error {
	extraFiles, err := createListenerFiles(s.Listeners)
	if err != nil {
		return err
	}
	defer closeListenerFiles(extraFiles)
	if err := s.startWorker(ctx, extraFiles); err != nil {
		return err
	}
	<-s.chanCloseMonitor.Done()
	return nil
}

// RestartProcess graceful restarts worker process
func (s *Supervisor) RestartProcess(ctx context.Context, stopSig os.Signal) error {
	if err := s.restartWorker(ctx, stopSig); err != nil {
		return err
	}
	return nil
}

// Shutdown worker process
func (s *Supervisor) Shutdown(ctx context.Context, stopSig os.Signal) error {
	if err := s.shutdownWorker(ctx, stopSig); err != nil {
		return err
	}
	return nil
}

func (s *Supervisor) startWorker(ctx context.Context, extraFiles []*os.File) error {
	s.workerMu.Lock() // worker LOCK
	wk := &worker.Worker{
		Command:            s.Command,
		Args:               s.Args,
		ExtraFiles:         extraFiles,
		WaitReadyFunc:      s.WaitReadyFunc,
		AutoRestartTimeout: s.AutoRestartTimeout,
	}
	wk.SetAutoRestart(s.AutoRestartEnabled)
	s.worker = wk
	s.workerMu.Unlock() // worker UNLOCK

	if err := wk.Start(ctx); err != nil {
		return fmt.Errorf("supervisor: failed to start new worker: %v", err)
	}
	s.chanCloseMonitor.addDone(wk.Done())
	return nil
}

func (s *Supervisor) restartWorker(ctx context.Context, stopSig os.Signal) error {
	// renew worker
	s.workerMu.Lock() // worker LOCK
	oldwk := s.worker
	newwk := &worker.Worker{
		Command:            s.Command,
		Args:               s.Args,
		ExtraFiles:         oldwk.ExtraFiles,
		WaitReadyFunc:      s.WaitReadyFunc,
		AutoRestartTimeout: s.AutoRestartTimeout,
	}
	newwk.SetAutoRestart(s.AutoRestartEnabled)
	s.worker = newwk
	s.workerMu.Unlock() // worker UNLOCK

	if err := newwk.Start(ctx); err != nil {
		return fmt.Errorf("supervisor: failed to start new worker: %v", err)
	}
	s.chanCloseMonitor.addDone(newwk.Done())
	// stop old worker
	time.Sleep(s.StopOldDelay)
	if err := oldwk.Stop(ctx, stopSig); err != nil {
		log.Printf("supervisor: failed to stop old worker. sig %s. %v", stopSig, err)
		log.Println("supervisor: force stopping old worker")
		if err := oldwk.Kill(); err != nil {
			return fmt.Errorf("supervisor: faield to kill old worker: %v", err)
		}
	}
	return nil
}

func (s *Supervisor) shutdownWorker(ctx context.Context, stopSig os.Signal) error {
	s.workerMu.Lock() // worker LOCK
	wk := s.worker
	s.workerMu.Unlock() // worker UNLOCK

	if err := wk.Stop(ctx, stopSig); err != nil {
		log.Printf("supervisor: failed to stop worker. sig %s. %v", stopSig, err)
		log.Println("supervisor: force stopping worker")
		if err := wk.Kill(); err != nil {
			return fmt.Errorf("supervisor: faield to kill worker: %v", err)
		}
	}
	return nil
}

func createListenerFiles(listeners []net.Listener) ([]*os.File, error) {
	fs := make([]*os.File, 0)
	var err error
loop:
	for _, l := range listeners {
		switch l := l.(type) {
		case *net.TCPListener:
			f, e := l.File()
			if e != nil {
				err = fmt.Errorf("supervisor: failed to create listener file: %v", e)
				break loop
			}
			fs = append(fs, f)
		default:
			err = fmt.Errorf("supervisor: failed to create listener file. not implemented %T", l)
			break loop
		}
	}
	if err != nil {
		for _, f := range fs {
			if err := f.Close(); err != nil {
				log.Printf("supervisor: an error occurred %v", err)
			}
		}
		return nil, err
	}
	return fs, nil
}

func closeListenerFiles(files []*os.File) {
	for _, f := range files {
		if err := f.Close(); err != nil {
			log.Printf("supervisor: failed to close listener files: %v", err)
		}
	}
}

type chanCloseMonitor struct {
	wg   sync.WaitGroup
	done chan struct{}
	mu   sync.Mutex
}

func (dm *chanCloseMonitor) addDone(ch <-chan struct{}) {
	dm.wg.Add(1)
	go func() {
		defer dm.wg.Done()
		for range ch {
		}
	}()
}

func (dm *chanCloseMonitor) Done() <-chan struct{} {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	if dm.done == nil {
		dm.done = make(chan struct{})
		go func() {
			dm.wg.Wait()
			close(dm.done)
		}()
	}
	return dm.done
}
