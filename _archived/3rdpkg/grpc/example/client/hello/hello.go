package main

import (
	"log"

	pb "github.com/kei2100/playground-go/3rdpkg/grpc/example/hello"

	"github.com/kei2100/playground-go/3rdpkg/grpc/example/client/internal"
	"golang.org/x/net/context"
)

func main() {
	conn, err := internal.Connection()
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "kei2100"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}
