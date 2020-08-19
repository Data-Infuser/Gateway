package main

import (
	"gitlab.com/promptech1/infuser-gateway/client"
	"gitlab.com/promptech1/infuser-gateway/handler"
	"gitlab.com/promptech1/infuser-gateway/router"
)

func main() {
	grpcPool := client.NewGRPCPool()

	r := router.New()

	apiGroup := r.Group("/api")
	h := handler.NewHandler(grpcPool)
	h.Register(apiGroup)

	r.Logger.Fatal(r.Start(":8080"))
}
