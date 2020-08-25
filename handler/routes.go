package handler

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/promptech1/infuser-gateway/router/middleware"
)

func (h *Handler) Register(apiGroup *echo.Group) {
	apiGroup.Use(middleware.KeyExtractor())

	apiGroup.GET("/:nameSpace/:operation", h.ExecuteApi)
	apiGroup.GET("/:nameSpace/:operation/meta", h.ExecuteApi)
}

