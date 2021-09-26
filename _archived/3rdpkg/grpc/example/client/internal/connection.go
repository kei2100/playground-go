package internal

import (
	"crypto/tls"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Connection creates gRPC client connection.
func Connection() (*grpc.ClientConn, error) {
	port := os.Getenv("GRPC_PORT")
	if len(port) == 0 {
		port = "50051"
	}

	host := os.Getenv("GRPC_HOST")
	if len(host) == 0 {
		host = "localhost"
	}

	return grpc.Dial(fmt.Sprintf("%s:%s", host, port), dialOption())
}

func dialOption() grpc.DialOption {
	if os.Getenv("GRPC_USE_TLS") != "true" {
		return grpc.WithInsecure()
	}

	conf := &tls.Config{}

	if os.Getenv("GRPC_USE_TLS_INSECURE") == "true" {
		conf.InsecureSkipVerify = true
	}

	creds := credentials.NewTLS(conf)
	return grpc.WithTransportCredentials(creds)
}
