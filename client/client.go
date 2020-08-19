package client

import (
	"flag"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
	"time"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("server_addr", "localhost:9090", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name used to verify the hostname returned by the TLS handshake")
)

func NewGRPCPool() *Pool{
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
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
		conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
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
