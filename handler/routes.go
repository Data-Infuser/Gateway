package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gitlab.com/promptech1/infuser-gateway/router/middleware"
	"net/http/httputil"
	"net/url"
)

// RegisterAuth: 공공데이터 포털 인증키 연계 처리
func (h *Handler) RegisterAuth(group *echo.Group) {
	group.Use(middleware.KeyExtractor(h.ctx))

	group.POST("", h.RegisterPortalKey)
}

func (h *Handler) Rp(c echo.Context, endPoint string) error {
	url, _ := url.Parse(endPoint)
	req := c.Request()
	res := c.Response()
	// create the reverse proxy
	proxyUrl, _ := url.Parse(fmt.Sprintf("%s://%s", url.Scheme, url.Host))
	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.RawQuery = c.QueryParams().Encode()
	req.URL.Path = url.Path
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
	return nil
}

// RegisterApi: Rest API 대응을 위한 route 정보 정의
func (h *Handler) RegisterApi(apiGroup *echo.Group) {
	apiGroup.Use(middleware.KeyExtractor(h.ctx))

	apiGroup.GET("/:nameSpace/:version/:operation", h.ExecuteAPI)
	apiGroup.GET("/:nameSpace/:version/:operation/meta", h.FindMeta)
}

