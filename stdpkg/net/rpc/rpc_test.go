package rpc

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"testing"
)

// TestService is a service of the RPC
type TestService struct{}

// Call is a service method of the TestService
func (s *TestService) Call(msg string, reply *string) error {
	*reply = fmt.Sprintf("reply: %s", msg)
	return nil
}

type server struct{}

func (s *server) serve(ln net.Listener) error {
	srv := rpc.NewServer()
	srv.Register(&TestService{})
	srv.Accept(ln)
	return nil
}

type client struct {
	cli *rpc.Client
}

func newClient(network, addr string) (*client, error) {
	cli, err := rpc.Dial(network, addr)
	if err != nil {
		return nil, fmt.Errorf("dial error: %v", err)
	}
	return &client{cli: cli}, nil
}

func (c *client) call(msg string) (string, error) {
	var reply string
	err := c.cli.Call("TestService.Call", "hello", &reply)
	if err != nil {
		return "", fmt.Errorf("call error: %v", err)
	}
	return reply, nil
}

func (c *client) close() error {
	return c.cli.Close()
}

func TestRPC(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		srv := server{}
		srv.serve(ln)
		log.Println("serve done")
	}()

	cli, err := newClient("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer cli.close()
	reply, err := cli.call("hello")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(reply)
	ln.Close()

	<-done
	log.Println("done")
}
