// handler: Rest API 요청 처리 및 gRPC 통신을 통한 데이터 교환 수행 package
package handler

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/promptech1/infuser-gateway/enum"
	grpc_author "gitlab.com/promptech1/infuser-gateway/infuser-protobuf/gen/proto/author"
)

// ExecuteApi: 활용자의 Rest API 요청을 처리함. gRPC를 통해 필요한 데이터를 교환하고 그결과를 JSON 형태로 반환함
func (h *Handler) ExecuteApi(c echo.Context) error {
	ctx := c.Request().Context()
	c.Logger().Debug("Test logger ================")

	token, _ := c.Get("Token").(string)

	conn, err := h.pool.Get(ctx)
	if err != nil {
		code := enum.InternalException
		return c.JSON(code.HttpCode(), map[string]interface{}{
			"code": code,
			"msg":  code.Message(),
		})
	}
	defer conn.Close()

	client := grpc_author.NewApiAuthServiceClient(conn)

	apiAuthRes, err := client.Auth(ctx, &grpc_author.ApiAuthReq{
		NameSpace:    c.Param("nameSpace"),
		OperationUrl: c.Param("operation"),
		Token:        token,
	})

	var code enum.ResCode
	if err != nil {
		code = enum.Unknown
	} else {
		code = enum.FindResCode(apiAuthRes.Code)
	}

	if !code.Valid() {
		return c.JSON(code.HttpCode(), map[string]interface{}{
			"code": code,
			"msg":  code.Message(),
		})
	} else {
		// TODO: Executer 연계를 통한 데이터 fetch 수행 필요함

		return c.JSON(code.HttpCode(), map[string]interface{}{
			"code": code,
			"msg":  code.Message(),
		})
	}

	return nil
}
