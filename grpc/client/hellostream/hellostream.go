package main

import (
	"log"

	pb "github.com/kei2100/playground-go/grpc/hello"

	"io"

	"github.com/kei2100/playground-go/grpc/client/internal"
	"golang.org/x/net/context"
)

func main() {
	conn, err := internal.Connection()
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	stream, err := c.SayHelloStream(context.Background())
	if err != nil {
		log.Fatalf("failed to create stream: %v", err)
	}

	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("failed to receive: %v", err)
			}
			log.Printf("Greeting: %s", in.Message)
		}
	}()

	for _, name := range []string{"kei2100", "2100kei"} {
		if err := stream.Send(&pb.HelloRequest{Name: name}); err != nil {
			log.Fatalf("failed to send: %v", err)
		}
	}
	if err := stream.CloseSend(); err != nil {
		log.Fatalf("failed to close send: %v", err)
	}
	<-waitc
}
