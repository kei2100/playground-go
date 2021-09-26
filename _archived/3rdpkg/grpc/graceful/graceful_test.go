package graceful

import (
	"net"
	"sync"
	"testing"
	"time"

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

func TestGracefulStopBehavior(t *testing.T) {
	s := grpc.NewServer()
	ss := &stubServer{emptyCalled: make(chan struct{}, 2), resumeEmptyCall: make(chan struct{})}
	testpb.RegisterTestServiceServer(s, ss)
	reflection.Register(s)

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	println(ln.Addr().String()) // for check tcp connections
	go s.Serve(ln)

	cc, err := grpc.Dial(ln.Addr().String(), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer cc.Close()
	c := testpb.NewTestServiceClient(cc)

	var wg sync.WaitGroup
	wg.Add(2)

	// GracefulStopする前のRPCコール
	go func() {
		defer wg.Done()
		if _, err := c.EmptyCall(context.Background(), new(testpb.Empty)); err != nil {
			// resumeEmptyCallした後は、エラー無しでレスポンスされること
			t.Errorf("got err %v, want no error", err)
		}
	}()
	select {
	case <-ss.emptyCalled:
		// ok
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout exceeded while waiting for receive emptyCalled")
	}

	// ブレークポイントで止めると、GracefulStop前は例えばこんな状態
	//
	// tcp4       0      0  127.0.0.1.51174        127.0.0.1.51175        ESTABLISHED
	// tcp4       0      0  127.0.0.1.51175        127.0.0.1.51174        ESTABLISHED
	// tcp4       0      0  127.0.0.1.51174        *.*                    LISTEN
	gracefulDone := make(chan struct{})
	go func() {
		// GracefulStopするとlistenerがクローズされる。
		// 保留中のRPCコールが完了するまでこのメソッドはブロックする。
		s.GracefulStop()
		close(gracefulDone)
	}()

	// GracefulStop後はLISTENが止まる
	//
	// tcp4       0      0  127.0.0.1.51174        127.0.0.1.51175        ESTABLISHED
	// tcp4       0      0  127.0.0.1.51175        127.0.0.1.51174        ESTABLISHED
	time.Sleep(10 * time.Millisecond) // GracefulStopのgoroutineが進捗してリスンが停止する挙動を見やすいようにsleep
	close(ss.resumeEmptyCall)

	select {
	case <-gracefulDone:
		// ok
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout exceeded while waiting for graceful stop")
	}

	// 保留していたRPCコールが完了して、GracefulStopが完了したら、既存のコネクションはクローズされる

	// GracefulStop後のRPCコールはエラーになること
	go func() {
		defer wg.Done()
		if _, err := c.EmptyCall(context.Background(), new(testpb.Empty)); err == nil {
			t.Error("got nil, want an error")
		}
	}()

	wg.Wait()
}
