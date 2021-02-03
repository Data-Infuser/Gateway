package handler

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/promptech1/infuser-gateway/router/middleware"
)

// RegisterAuth: 공공데이터 포털 인증키 연계 처리
func (h *Handler) RegisterAuth(group *echo.Group) {
	group.Use(middleware.KeyExtractor(h.ctx))

	group.POST("", h.RegisterPortalKey)
}

// RegisterApi: Rest API 대응을 위한 route 정보 정의
func (h *Handler) RegisterApi(apiGroup *echo.Group) {
	apiGroup.Use(middleware.KeyExtractor(h.ctx))

	apiGroup.GET("/:nameSpace/:version/:operation", h.ExecuteAPI)
	apiGroup.GET("/:nameSpace/:version/:operation/meta", h.FindMeta)
}
