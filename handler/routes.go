package handler

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/promptech1/infuser-gateway/router/middleware"
)

func (h *Handler) Register(apiGroup *echo.Group) {
	apiGroup.Use(middleware.KeyExtractor())

	apiGroup.GET("/:nameSpace/:service", h.ExecuteApi)
	apiGroup.GET("/:nameSpace/:service/meta", h.ExecuteApi)
}

