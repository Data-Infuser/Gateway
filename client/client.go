package client

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"gitlab.com/promptech1/infuser-gateway/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
	"time"
)

var (
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name used to verify the hostname returned by the TLS handshake")
)

func NewGRPCPool(conf *config.Config) *Pool {
	var opts []grpc.DialOption
	if conf.Author.Tls {
		if conf.Author.CaFile == "" {
			conf.Author.CaFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(conf.Author.CaFile, *serverHostOverride)
		if err != nil {
			glog.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithBlock())

	var factory Factory
	factory = func() (*grpc.ClientConn, error) {
		serverAddr := fmt.Sprintf("%s:%d", conf.Author.Host, conf.Author.Port)
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		if err != nil {
			glog.Infoln("Failed to start gRPC connection: %v", err)
		}
		return conn, err
	}

	pool, err := New(factory, 10, 50, time.Second)
	if err != nil {
		glog.Infoln("Failed to create gRPC pool: %v", err)
	}

	return pool
}
