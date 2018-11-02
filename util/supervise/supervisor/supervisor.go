package supervisor

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

type supervisor struct {
	exec string
	argv []string

	worker   *exec.Cmd
	workerMu sync.RWMutex

	stop   bool
	stopMu sync.RWMutex
}

// NewSupervisor creates a supervisor
func NewSupervisor(exec string, argv ...string) *supervisor {
	return &supervisor{
		exec: exec,
		argv: argv,
	}
}

// Start supervisor
func (sv *supervisor) Start() error {
	log.Printf("supervisor: started. pid %d", os.Getpid())
	if err := sv.startWorker(); err != nil {
		return err
	}
	go sv.handleSignal()
	for {
		err := sv.waitWorker()
		if err != nil {
			log.Printf("supervisor: woker abnormally finished. %v", err)
		}
		if sv.stopCalled() {
			return nil
		}
		log.Println("supervisor: restarting the worker...")
		if err := sv.startWorker(); err != nil {
			return err
		}
	}
}

func (sv *supervisor) startWorker() error {
	sv.workerMu.Lock()
	defer sv.workerMu.Unlock()
	c := exec.Command(sv.exec, sv.argv...)
	c.Stdout, c.Stderr = os.Stdout, os.Stderr
	if err := c.Start(); err != nil {
		return fmt.Errorf("supervisor: failed to start the worker: %v", err)
	}
	log.Printf("supervisor: worker started. pid %d", c.Process.Pid)
	sv.worker = c
	return nil
}

func (sv *supervisor) waitWorker() error {
	sv.workerMu.RLock()
	w := sv.worker
	sv.workerMu.RUnlock()
	return w.Wait()
}

func (sv *supervisor) signalWorker(sig os.Signal) error {
	sv.workerMu.RLock()
	defer sv.workerMu.RUnlock()
	return sv.worker.Process.Signal(sig)
}

func (sv *supervisor) handleSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	for sig := range ch {
		if sig == syscall.SIGHUP {
			// do something...
			continue
		}
		sv.stopCall()
		log.Println("supervisor: stopping the supervisor")
		if err := sv.signalWorker(sig); err != nil {
			log.Printf("supervisor: failed to send sig %s to the worker: %v", sig, err)
		}
		return
	}
}

func (sv *supervisor) stopCall() {
	sv.stopMu.Lock()
	defer sv.stopMu.Unlock()
	sv.stop = true
}

func (sv *supervisor) stopCalled() bool {
	sv.stopMu.RLock()
	defer sv.stopMu.RUnlock()
	return sv.stop
}
