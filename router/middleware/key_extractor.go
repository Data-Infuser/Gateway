// middleware: handler 모듈 구동에 필요한 middleware 정의
package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/promptech1/infuser-gateway/config"
	"gitlab.com/promptech1/infuser-gateway/enum"
)

type (
	keyExtractor func(echo.Context) (string, error)
)

var (
	ErrTokenMissing = echo.NewHTTPError(http.StatusUnauthorized, "missing service key")
)

// KeyExtractor ... : Rest API를 통해 전송된(Request Header 또는 Query Param) 사용자 인증키를 추출
func KeyExtractor(ctx *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		headerExtractor := keyFromHeader("Authorization", "Infuser")
		querytExtractor := keyFromQuery("ServiceKey")

		return func(c echo.Context) error {
			var token string

			if ctx.Server.ExceptServiceKey {
				c.Set("Token", ctx.Author.MasterKey)
				return next(c)
			}

			token, _ = headerExtractor(c)
			if len(token) > 0 {
				c.Set("Token", token)
				return next(c)
			}

			token, _ = querytExtractor(c)
			if len(token) > 0 {
				c.Set("Token", token)
				return next(c)
			}

			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code": enum.Unauthorized,
				"msg":  "인증키는 필수 항목 입니다.",
			})
		}
	}
}

// keyFromHeader: request header로 부터 인증키 추출
func keyFromHeader(header string, authScheme string) keyExtractor {
	return func(c echo.Context) (string, error) {
		// Header로 부터 api key 추출
		auth := c.Request().Header.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}

		return "", ErrTokenMissing
	}
}

// keyFromQuery: query parameter로 부터 인증키 추출
func keyFromQuery(authScheme string) keyExtractor {
	return func(c echo.Context) (string, error) {
		auth := c.QueryParam(authScheme)
		if len(auth) > 0 {
			return auth, nil
		}

		return "", ErrTokenMissing
	}
}
