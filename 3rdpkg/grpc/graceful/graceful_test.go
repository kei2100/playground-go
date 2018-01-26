package graceful

import (
	"net"
	"testing"

	"sync"
	"time"

	"github.com/kei2100/playground-go/util/wait"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	testpb "google.golang.org/grpc/test/grpc_testing"
)

type stubServer struct {
	testpb.TestServiceServer

	emptyCalled     chan struct{}
	resumeEmptyCall chan struct{}
}

func (s *stubServer) EmptyCall(context.Context, *testpb.Empty) (*testpb.Empty, error) {
	s.emptyCalled <- struct{}{}
	<-s.resumeEmptyCall
	return new(testpb.Empty), nil
}

func TestTest(t *testing.T) {
	s := grpc.NewServer()
	ss := &stubServer{emptyCalled: make(chan struct{}, 2), resumeEmptyCall: make(chan struct{})}
	testpb.RegisterTestServiceServer(s, ss)
	reflection.Register(s)

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	println(ln.Addr().String())
	go s.Serve(ln)

	cc, err := grpc.Dial(ln.Addr().String(), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer cc.Close()
	c := testpb.NewTestServiceClient(cc)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// GracefulStopする前のRPCコールは最後まで実行されること
		if _, err := c.EmptyCall(context.Background(), new(testpb.Empty)); err != nil {
			t.Errorf("got err %v, want no error", err)
		}
	}()
	if err := wait.ReceiveStruct(ss.emptyCalled, 100*time.Millisecond); err != nil {
		t.Fatal(err)
	}

	// GracefulStop前は例えばこんな状態
	//
	// tcp4       0      0  127.0.0.1.51174        127.0.0.1.51175        ESTABLISHED
	// tcp4       0      0  127.0.0.1.51175        127.0.0.1.51174        ESTABLISHED
	// tcp4       0      0  127.0.0.1.51174        *.*                    LISTEN
	go s.GracefulStop()

	// GracefulStop後はLISTENが止まる
	//
	// tcp4       0      0  127.0.0.1.51174        127.0.0.1.51175        ESTABLISHED
	// tcp4       0      0  127.0.0.1.51175        127.0.0.1.51174        ESTABLISHED
	time.Sleep(10 * time.Millisecond)

	close(ss.resumeEmptyCall)

	// 保留していたRPCコールが完了してもFINやRSTでコネクションクローズするわけではない。
	// 新規の接続や、RPCコールはエラーになる。
	//
	// tcp4       0      0  127.0.0.1.51174        127.0.0.1.51175        ESTABLISHED
	// tcp4       0      0  127.0.0.1.51175        127.0.0.1.51174        ESTABLISHED

	go func() {
		defer wg.Done()
		// GracefulStop後のRPCコールはエラーになること
		if _, err := c.EmptyCall(context.Background(), new(testpb.Empty)); err == nil {
			t.Error("got nil, want an error")
		}
	}()

	if err := wait.WGroup(&wg, 100*time.Millisecond); err != nil {
		t.Fatal(err)
	}
}
