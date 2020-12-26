package handler

import (
	"gitlab.com/promptech1/infuser-gateway/client"
	"gitlab.com/promptech1/infuser-gateway/config"
)

type Handler struct {
	authPool     *client.Pool
	executorPool *client.Pool
	ctx          *config.Config
}

// NewHandler: gRPC Pool을 매개변수로 한 Handler 객체 생성
func NewHandler(authPool *client.Pool, executorPool *client.Pool, ctx *config.Config) *Handler {
	return &Handler{authPool: authPool, executorPool: executorPool, ctx: ctx}
}
