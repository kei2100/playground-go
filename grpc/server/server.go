package main

import (
	"log"
	"net"

	pb "github.com/kei2100/playground-go/grpc/hello"

	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
)

// server is used to implement hello.GreeterServer.
type server struct{}

// SayHello implements hello.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	port := os.Getenv("GRPC_PORT")
	if len(port) == 0 {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":" + port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
