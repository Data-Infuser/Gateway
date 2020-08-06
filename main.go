package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"gitlab.com/promptech1/infuser-gateway/client"
	"gitlab.com/promptech1/infuser-gateway/enum"
	grpc_author "gitlab.com/promptech1/infuser-gateway/infuser-protobuf/gen/proto/author"
	"strings"
)

func main() {
	ctx := context.Background()

	grpcPool := client.NewGRPCPool()

	router := gin.Default()
	router.GET("/rest/:nameSpace/:service", func(c *gin.Context) {
		var hasToken = true

		// Token값은 Request Header 또는 Query Parameter 형태로 전송
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			if token = c.Query("token"); token == "" {
				hasToken = false
			}
		} else {
			i := strings.Index(token, "Bearer ")
			if i == 0 {
				token = strings.Split(token, "Bearer ")[1]
			} else {
				hasToken = false
			}
		}

		// Token 값이 없는 경우 오류 처리
		if !hasToken {
			code := enum.Unauthorized
			c.JSON(code.HttpCode(), gin.H {
				"code": code,
				"msg": "인증키는 필수 항목 입니다.",
			})
			return
		}

		// Grpc 통신을 통한 인증 처리
		conn, err := grpcPool.Get(ctx)
		_ = err
		defer conn.Close()

		client := grpc_author.NewApiAuthServiceClient(conn)

		apiAuthRes, err := client.Auth(ctx, &grpc_author.ApiAuthReq{
			NameSpace: c.Param("nameSpace"),
			Token: token,
		})

		var code enum.ResCode
		if err != nil {
			code = enum.Unknown
		} else {
			code = enum.FindResCode(apiAuthRes.Code)
		}

		if !code.Valid() {
			c.JSON(code.HttpCode(), gin.H {
				"code": code,
				"msg": code.Message(),
			})
		} else {
			// TODO: Executer 연계를 통한 데이터 fetch 수행 필요함

			c.JSON(code.HttpCode(), gin.H {
				"code": code,
				"msg": code.Message(),
			})
		}
	})

	router.Run(":8080")
}
