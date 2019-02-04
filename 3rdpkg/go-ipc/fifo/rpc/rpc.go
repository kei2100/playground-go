package main

import (
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"
	"os/exec"
	"sync"

	"bitbucket.org/avd/go-ipc/fifo"
)

// rpc over fifo
func main() {
	if os.Getenv("_WORKER") == "on" {
		log.Printf("pid %v: start worker", os.Getpid())

		conn, err := newFIFOConn("rfifo", "wfifo")
		panicIf(err)
		defer conn.Close()

		log.Printf("pid %v: serveRPC", os.Getpid())
		serveRPC(conn)
		log.Printf("pid %v: stop serveRPC", os.Getpid())
		return
	}

	log.Printf("pid %v: start master", os.Getpid())

	wkp, err := launchWorker()
	panicIf(err)
	defer func() {
		wkp.Wait()
	}()

	conn, err := newFIFOConn("wfifo", "rfifo") // swap r, w
	panicIf(err)
	defer conn.Close()

	log.Printf("pid %v: callRPC", os.Getpid())
	reply, err := callRPC(conn)
	log.Printf("pid %v: callRPC reply %v", os.Getpid(), reply)
}

func newFIFOConn(rFIFO, wFIFO string) (*fifoConn, error) {
	var rfifo, wfifo fifo.Fifo
	var rerr, werr error

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		rfifo, rerr = fifo.New(rFIFO, os.O_CREATE|os.O_RDONLY, 0666)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		wfifo, werr = fifo.New(wFIFO, os.O_CREATE|os.O_WRONLY, 0666)
	}()
	wg.Wait()

	if rerr != nil || werr != nil {
		if rfifo != nil {
			rfifo.Close()
		}
		if wfifo != nil {
			wfifo.Close()
		}
		return nil, fmt.Errorf("rerr: %v, werr: %v", rerr, werr)
	}
	return &fifoConn{rfifo: rfifo, wfifo: wfifo}, nil
}

type fifoConn struct {
	rfifo fifo.Fifo
	wfifo fifo.Fifo
}

func (conn *fifoConn) Read(p []byte) (int, error) {
	return conn.rfifo.Read(p)
}

func (conn *fifoConn) Write(p []byte) (int, error) {
	return conn.wfifo.Write(p)
}

func (conn *fifoConn) Close() error {
	// wからcloseすること。rからだとwinでハングするため
	werr := conn.wfifo.Close()
	rerr := conn.rfifo.Close()

	if rerr != nil || werr != nil {
		return fmt.Errorf("close error: rerr %v, werr %v", rerr, werr)
	}
	return nil
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func launchWorker() (*os.Process, error) {
	// launch worker
	bin, err := os.Executable()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(bin)
	cmd.Env = append(cmd.Env, "_WORKER=on")
	cmd.Env = append(cmd.Env, "TMPDIR="+os.TempDir())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd.Process, nil
}

// MessageService is a RPC Service
type MessageService struct{}

// Echo RPC method
func (s *MessageService) Echo(arg *EchoArg, reply *EchoReply) error {
	reply.Message = fmt.Sprintf("%s from worker", arg.Message)
	return nil
}

// EchoArg is an argument for the MessageService.Echo
type EchoArg struct {
	Message string
}

// EchoReply is a reply for the MessageService.Echo
type EchoReply struct {
	Message string
}

func callRPC(conn io.ReadWriteCloser) (string, error) {
	cli := rpc.NewClient(conn)
	defer cli.Close()
	arg, reply := EchoArg{}, EchoReply{}
	arg.Message = "hello"
	if err := cli.Call("MessageService.Echo", &arg, &reply); err != nil {
		return "", fmt.Errorf("m: failed to call: %v", err)
	}
	return reply.Message, nil
}

func serveRPC(conn io.ReadWriteCloser) error {
	srv := rpc.NewServer()
	srv.Register(&MessageService{})
	srv.ServeConn(conn)
	return nil
}
