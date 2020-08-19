package handler

import (

	"gitlab.com/promptech1/infuser-gateway/client"
)

type Handler struct {
	pool *client.Pool
}

func NewHandler(pool *client.Pool) *Handler {
	return &Handler{pool: pool}
}
