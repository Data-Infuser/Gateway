package client

import (
	"flag"
	"fmt"
	"time"

	"gitlab.com/promptech1/infuser-gateway/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

var (
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name used to verify the hostname returned by the TLS handshake")
)

func NewGRPCPool(ctx *config.Context) *Pool {
	var opts []grpc.DialOption
	if ctx.Author.Tls {
		if ctx.Author.CaFile == "" {
			ctx.Author.CaFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(ctx.Author.CaFile, *serverHostOverride)
		if err != nil {
			ctx.Logger.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithBlock())

	var factory Factory
	factory = func() (*grpc.ClientConn, error) {
		serverAddr := fmt.Sprintf("%s:%d", ctx.Author.Host, ctx.Author.Port)
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		if err != nil {
			ctx.Logger.Fatalf("Failed to start gRPC connection: %v", err)
		}
		return conn, err
	}

	pool, err := New(factory, 10, 50, time.Second)
	if err != nil {
		ctx.Logger.Fatalf("Failed to create gRPC pool: %v", err)
	}

	return pool
}
