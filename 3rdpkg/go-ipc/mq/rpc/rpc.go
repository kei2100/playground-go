package main

import (
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"
	"os/exec"

	"bitbucket.org/avd/go-ipc/mq"
)

// MessageService is a RPC Service
type MessageService struct{}

// Echo RPC method
func (s *MessageService) Echo(arg *EchoArg, reply *EchoReply) error {
	reply.Message = fmt.Sprintf("%s from %d", arg.Message, os.Getpid())
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

// Master represents a parent process
type Master struct{}

// Worker represents a child process
type Worker struct{}

// Exec master
func (m *Master) Exec(cli *rpc.Client) error {
	log.Println("m: call")
	arg, reply := EchoArg{}, EchoReply{}
	arg.Message = "hello"
	if err := cli.Call("MessageService.Echo", &arg, &reply); err != nil {
		return fmt.Errorf("m: failed to call: %v", err)
	}
	log.Printf("m: reply %s", reply.Message)
	return nil
}

// Exec worker
func (w *Worker) Exec(conn io.ReadWriteCloser) error {
	srv := rpc.NewServer()
	srv.Register(&MessageService{})
	log.Println("w: serve")
	srv.ServeConn(conn)
	log.Println("w: serve stopped")
	return nil
}

// mqConn bridges between mq.Messenger and io.ReadWriteCloser
type mqConn struct {
	rmq mq.Messenger
	wmq mq.Messenger
}

func (conn *mqConn) Read(p []byte) (int, error) {
	return conn.rmq.Receive(p)
}

func (conn *mqConn) Write(p []byte) (int, error) {
	err := conn.wmq.Send(p)
	return len(p), err
}

func (conn *mqConn) Close() error {
	conn.rmq.Close()
	conn.wmq.Close()
	// FIX
	return nil
}

func main() {
	if os.Getenv("_WORKER") == "on" {
		log.Printf("w: pid %v", os.Getpid())
		// -- worker process --
		wmq, err := mq.Open("rmq", 0)
		if err != nil {
			log.Fatalf("w: failed to open wmq: %v", err)
		}
		rmq, err := mq.Open("wmq", 0)
		if err != nil {
			log.Fatalf("w: failed to open rmq: %v", err)
		}
		conn := &mqConn{rmq: rmq, wmq: wmq}
		defer conn.Close()

		w := Worker{}
		if err := w.Exec(conn); err != nil {
			log.Fatal(err)
		}
		log.Println("w: done exec")
		return
	}

	// -- master process --
	log.Printf("m: pid %v", os.Getpid())
	mq.Destroy("rmq")
	mq.Destroy("wmq")
	rmq, err := mq.New("rmq", os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		log.Fatalf("m: failed to new rmq: %v", err)
	}
	wmq, err := mq.New("wmq", os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		log.Fatalf("m: failed to new wmq: %v", err)
	}

	conn := &mqConn{rmq: rmq, wmq: wmq}
	defer conn.Close()

	cli := rpc.NewClient(conn)
	defer cli.Close()

	// launch worker
	bin, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(bin)
	cmd.Env = append(cmd.Env, "_WORKER=on")
	cmd.Env = append(cmd.Env, "TMPDIR="+os.TempDir())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		//cmd.Process.Signal(os.Interrupt)
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}
	}()

	master := &Master{}
	if err := master.Exec(cli); err != nil {
		cmd.Process.Kill()
		log.Fatal(err)
	}

	log.Println("m: done exec")
}
