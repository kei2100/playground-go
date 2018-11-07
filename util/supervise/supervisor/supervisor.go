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

// Supervisor struct
type Supervisor struct {
	exec string
	args []string

	worker   *exec.Cmd
	workerMu sync.RWMutex

	stop   bool
	stopMu sync.RWMutex
}

// NewSupervisor creates a Supervisor
func NewSupervisor(exec string, args ...string) *Supervisor {
	return &Supervisor{
		exec: exec,
		args: args,
	}
}

// Start Supervisor
func (sv *Supervisor) Start() error {
	log.Printf("Supervisor: started. pid %d", os.Getpid())
	if err := sv.startWorker(); err != nil {
		return err
	}
	go sv.handleSignal()
	for {
		err := sv.waitWorker()
		if err != nil {
			log.Printf("Supervisor: woker abnormally finished. %v", err)
		}
		if sv.stopCalled() {
			return nil
		}
		log.Println("Supervisor: restarting the worker...")
		if err := sv.startWorker(); err != nil {
			return err
		}
	}
}

func (sv *Supervisor) startWorker() error {
	sv.workerMu.Lock()
	defer sv.workerMu.Unlock()
	c := exec.Command(sv.exec, sv.args...)
	c.Stdout, c.Stderr = os.Stdout, os.Stderr
	if err := c.Start(); err != nil {
		return fmt.Errorf("Supervisor: failed to start the worker: %v", err)
	}
	log.Printf("Supervisor: worker started. pid %d", c.Process.Pid)
	sv.worker = c
	return nil
}

func (sv *Supervisor) waitWorker() error {
	sv.workerMu.RLock()
	w := sv.worker
	sv.workerMu.RUnlock()
	return w.Wait()
}

func (sv *Supervisor) signalWorker(sig os.Signal) error {
	sv.workerMu.RLock()
	defer sv.workerMu.RUnlock()
	return sv.worker.Process.Signal(sig)
}

func (sv *Supervisor) handleSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	for sig := range ch {
		if sig == syscall.SIGHUP {
			// do something...
			continue
		}
		sv.stopCall()
		log.Println("Supervisor: stopping the Supervisor")
		if err := sv.signalWorker(sig); err != nil {
			log.Printf("Supervisor: failed to send sig %s to the worker: %v", sig, err)
		}
		return
	}
}

func (sv *Supervisor) stopCall() {
	sv.stopMu.Lock()
	defer sv.stopMu.Unlock()
	sv.stop = true
}

func (sv *Supervisor) stopCalled() bool {
	sv.stopMu.RLock()
	defer sv.stopMu.RUnlock()
	return sv.stop
}
