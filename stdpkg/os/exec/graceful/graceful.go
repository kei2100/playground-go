package graceful

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kei2100/playground-go/stdpkg/os/exec/graceful/supervisor"
)

// Serve executes given command
// and graceful restarts when the restart signal received.
// default restart signal is HUP.
func Serve(command string, opts ...OptionFunc) error {
	return graceful.Serve(command, opts...)
}

// Restart graceful restarts manually
func Restart() error {
	return graceful.Restart()
}

var graceful = NewGraceful()

// Graceful restart engine
type Graceful struct {
	manualRestartCh   chan struct{}
	manualRestartedCh chan error
}

// NewGraceful creates a new Graceful
func NewGraceful() *Graceful {
	return &Graceful{
		manualRestartCh:   make(chan struct{}),
		manualRestartedCh: make(chan error),
	}
}

// Serve executes given command
// and graceful restarts when the restart signal received.
// default restart signal is HUP.
func (g *Graceful) Serve(command string, opts ...OptionFunc) error {
	o := &option{}
	o.applyOrDefault(opts)

	extraFiles, err := createListenerFiles(o.listeners)
	if err != nil {
		return err
	}
	defer closeListenerFiles(extraFiles)

	sv := &supervisor.Supervisor{
		Command:            command,
		Args:               o.args,
		ExtraFiles:         extraFiles,
		Env:                []string{listenersEnv(o.listeners)},
		WaitReadyFunc:      o.waitReadyFunc,
		AutoRestartEnabled: o.autoRestartEnabled,
		AutoRestartTimeout: o.autoRestartTimeout,
		StopOldDelay:       o.stopOldDelay,
	}
	done := make(chan struct{})
	go func() {
		if err := start(sv, o); err != nil {
			log.Println(err)
		}
		close(done)
	}()

	restartCh := make(chan os.Signal)
	signal.Notify(restartCh, o.restartSignals...)
	shutdownCh := make(chan os.Signal)
	signal.Notify(shutdownCh, o.shutdownSignals...)

	for {
		select {
		case <-done:
			return nil
		case <-restartCh:
			if err := restart(sv, o); err != nil {
				return err
			}
		case <-g.manualRestartCh:
			err := restart(sv, o)
			g.manualRestartedCh <- err
		case sig := <-shutdownCh:
			if err := shutdown(sv, sig, o); err != nil {
				return err
			}
			return nil
		}
	}
}

// Restart graceful restarts manually
func (g *Graceful) Restart() error {
	g.manualRestartCh <- struct{}{}
	return <-g.manualRestartedCh
}

func start(sv *supervisor.Supervisor, o *option) error {
	ctx, can := o.startContext()
	defer can()
	err := sv.Start(ctx)
	if err != nil {
		return fmt.Errorf("supervisor: failed to start process: %v", err)
	}
	return nil
}

func restart(sv *supervisor.Supervisor, o *option) error {
	ctx, can := o.restartContext()
	defer can()
	err := sv.RestartProcess(ctx, o.gracefulStopSignal)
	if err != nil {
		return fmt.Errorf("graceful: failed to restart process: %v", err)
	}
	return nil
}

func shutdown(sv *supervisor.Supervisor, sig os.Signal, o *option) error {
	ctx, can := o.shutdownContext()
	defer can()
	err := sv.Shutdown(ctx, sig)
	if err != nil {
		return fmt.Errorf("graceful: failed to shutdown process")
	}
	return nil
}

// options
type option struct {
	args          []string
	listeners     []net.Listener
	waitReadyFunc func(ctx context.Context, extraFileConns []net.Conn) error

	autoRestartEnabled bool
	autoRestartTimeout time.Duration

	restartSignals     []os.Signal
	shutdownSignals    []os.Signal
	gracefulStopSignal os.Signal

	startTimeout    time.Duration
	restartTimeout  time.Duration
	shutdownTimeout time.Duration
	stopOldDelay    time.Duration
}

func (o *option) applyOrDefault(opts []OptionFunc) {
	o.restartSignals = []os.Signal{syscall.SIGHUP}
	o.shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT}
	o.gracefulStopSignal = syscall.SIGTERM
	o.stopOldDelay = time.Second
	for _, f := range opts {
		f(o)
	}
}

var nopCancelFunc context.CancelFunc = func() {}

func (o *option) startContext() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	can := nopCancelFunc
	if o.startTimeout > 0 {
		ctx, can = context.WithTimeout(ctx, o.startTimeout)
	}
	return ctx, can
}

func (o *option) restartContext() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	can := nopCancelFunc
	if o.restartTimeout > 0 {
		ctx, can = context.WithTimeout(ctx, o.restartTimeout)
	}
	return ctx, can
}

func (o *option) shutdownContext() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	can := nopCancelFunc
	if o.shutdownTimeout > 0 {
		ctx, can = context.WithTimeout(ctx, o.shutdownTimeout)
	}
	return ctx, can
}

// OptionFunc is optional function for graceful
type OptionFunc func(o *option)

// WithArgs set command line arguments
func WithArgs(args ...string) OptionFunc {
	return func(o *option) { o.args = args }
}

// WithListeners set listeners.
// listeners are copied to os.File and set to extra files of worker process.
func WithListeners(listeners ...net.Listener) OptionFunc {
	return func(o *option) { o.listeners = listeners }
}

// WithWaitReadyFunc set WaitReadyFunc
func WithWaitReadyFunc(waitReadyFunc func(context.Context, []net.Conn) error) OptionFunc {
	return func(o *option) { o.waitReadyFunc = waitReadyFunc }
}

// WithAutoRestartEnabled set autoRestartEnabled
func WithAutoRestartEnabled(autoRestartEnabled bool) OptionFunc {
	return func(o *option) { o.autoRestartEnabled = autoRestartEnabled }
}

// WithAutoRestartTimeout set autoRestartTimeout
func WithAutoRestartTimeout(autoRestartTimeout time.Duration) OptionFunc {
	return func(o *option) { o.autoRestartTimeout = autoRestartTimeout }
}

// WithStartTimeout set startTimeout
func WithStartTimeout(startTimeout time.Duration) OptionFunc {
	return func(o *option) { o.startTimeout = startTimeout }
}

// WithRestartTimeout set restartTimeout
func WithRestartTimeout(restartTimeout time.Duration) OptionFunc {
	return func(o *option) { o.restartTimeout = restartTimeout }
}

// WithStopOldDelay set stopOldDelay
func WithStopOldDelay(stopOldDelay time.Duration) OptionFunc {
	return func(o *option) { o.stopOldDelay = stopOldDelay }
}
