package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"gitlab.com/promptech1/infuser-gateway/client"
	grpc_author "gitlab.com/promptech1/infuser-gateway/infuser-protobuf/gen/proto/author"
	"log"
)

func main() {
	ctx := context.Background()

	grpcPool := client.NewGRPCPool()

	router := gin.Default()
	router.GET("/rest/:nameSpace/:service", func(c *gin.Context) {
		log.Printf("Query %+v", c.Request.URL.Query())

		conn, err := grpcPool.Get(ctx)
		_ = err
		defer conn.Close()

		client := grpc_author.NewApiAuthServiceClient(conn)

		apiAuthRes, err := client.Auth(ctx, &grpc_author.ApiAuthReq{
			NameSpace: c.Param("nameSpace"),
			Token: c.Query("token"),
		})

		glog.Infof("gRPC res: +v", &apiAuthRes)

		c.JSON(200, gin.H{

		})
	})

	router.Run(":8080")
}
