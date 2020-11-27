package handler

import (
	"gitlab.com/promptech1/infuser-gateway/client"
)

type Handler struct {
	pool *client.Pool
}

// NewHandler: gRPC Pool을 매개변수로 한 Handler 객체 생성
func NewHandler(pool *client.Pool) *Handler {
	return &Handler{pool: pool}
}
