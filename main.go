package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
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

	r.HEAD("/hc", func(c echo.Context) error {
		return c.String(http.StatusOK, "OpenAPI Data Service")
	})

	h := handler.NewHandler(grpcAuthorPool, grpcExecutorPool, conf)
	authGroup := r.Group("/auth")
	h.RegisterAuth(authGroup)

	apiGroup := r.Group("/api")
	h.RegisterApi(apiGroup)

	r.Logger.Fatal(r.Start(conf.Server.Host + ":" + conf.Server.Port))
}
