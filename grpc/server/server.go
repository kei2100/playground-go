package main

import (
	"log"
	"net"

	pb "github.com/kei2100/playground-go/grpc/hello"

	"os"

	"fmt"

	"io"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// server is used to implement hello.GreeterServer.
type server struct{}

// SayHello implements hello.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func (s *server) SayHelloStream(stream pb.Greeter_SayHelloStreamServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		for _, str := range []string{"Good morning", "Good afternoon", "Good evening"} {
			if err := stream.Send(&pb.HelloReply{Message: fmt.Sprintf("%s %s", str, in.Name)}); err != nil {
				return err
			}
		}
	}
}

func main() {
	port := os.Getenv("GRPC_PORT")
	if len(port) == 0 {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
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
