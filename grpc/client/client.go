package main

import (
	"log"

	pb "github.com/kei2100/playground-go/grpc/hello"

	"os"

	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"crypto/tls"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if len(port) == 0 {
		port = "50051"
	}

	host := os.Getenv("GRPC_HOST")
	if len(host) == 0 {
		host = "localhost"
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), dialOption())
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

func dialOption() grpc.DialOption {
	if os.Getenv("GRPC_USE_TLS") != "true" {
		return grpc.WithInsecure()
	}

	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})
	return grpc.WithTransportCredentials(creds)
}
