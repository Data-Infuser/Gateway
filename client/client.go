// client:  Author 서버와의 통신을 위한 gRPC Client 정의
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

// NewGRPCAuthorPool: Config에 정의된 정보를 바탕으로 gRPC client을 생성하고 이를 관리할 수 있는 Pool을 반환함
func NewGRPCAuthorPool(ctx *config.Config) *Pool {
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

// NewGRPCAuthorPool: Config에 정의된 정보를 바탕으로 gRPC client을 생성하고 이를 관리할 수 있는 Pool을 반환함
func NewGRPCExecutorPool(ctx *config.Config) *Pool {
	var opts []grpc.DialOption
	if ctx.Executor.Tls {
		if ctx.Executor.CaFile == "" {
			ctx.Executor.CaFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(ctx.Executor.CaFile, *serverHostOverride)
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
		serverAddr := fmt.Sprintf("%s:%d", ctx.Executor.Host, ctx.Executor.Port)
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
