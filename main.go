package main

import (
	"log"
	"os"

	"gitlab.com/promptech1/infuser-gateway/client"
	"gitlab.com/promptech1/infuser-gateway/config"
	"gitlab.com/promptech1/infuser-gateway/handler"
	"gitlab.com/promptech1/infuser-gateway/router"
)

func main() {
	ballast := make([]byte, 10<<24)
	_ = ballast

	conf := new(config.Config)
	if err := conf.InitConf(); err != nil {
		log.Printf("Fail load config: %s", err.Error())
		os.Exit(-1)
	}

	grpcAuthorPool := client.NewGRPCAuthorPool(conf)
	grpcExecutorPool := client.NewGRPCExecutorPool(conf)

	r := router.New()

	apiGroup := r.Group("/api")
	h := handler.NewHandler(grpcAuthorPool, grpcExecutorPool, conf)
	h.Register(apiGroup)

	r.Logger.Fatal(r.Start(conf.Server.Host + ":" + conf.Server.Port))
}
