// handler: Rest API 요청 처리 및 gRPC 통신을 통한 데이터 교환 수행 package
package handler

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
	"gitlab.com/promptech1/infuser-gateway/enum"
	grpc_author "gitlab.com/promptech1/infuser-gateway/infuser-protobuf/gen/proto/author"
	grpc_executor "gitlab.com/promptech1/infuser-gateway/infuser-protobuf/gen/proto/executor"
)

// ExecuteApi: 활용자의 Rest API 요청을 처리함. gRPC를 통해 필요한 데이터를 교환하고 그결과를 JSON 형태로 반환함
func (h *Handler) ExecuteApi(c echo.Context) error {
	ctx := c.Request().Context()
	c.Logger().Debug("Test logger ================")

	token, _ := c.Get("Token").(string)

	authConn, err := h.authPool.Get(ctx)
	if err != nil {
		code := enum.InternalException
		return c.JSON(code.HttpCode(), map[string]interface{}{
			"code": code,
			"msg":  code.Message(),
		})
	}
	defer authConn.Close()

	executorConn, err := h.executorPool.Get(ctx)
	if err != nil {
		code := enum.InternalException
		return c.JSON(code.HttpCode(), map[string]interface{}{
			"code": code,
			"msg":  code.Message(),
		})
	}
	defer executorConn.Close()

	authorClient := grpc_author.NewApiAuthServiceClient(authConn)
	executorClient := grpc_executor.NewApiResultServiceClient(executorConn)

	apiAuthRes, err := authorClient.Auth(ctx, &grpc_author.ApiAuthReq{
		NameSpace:    c.Param("nameSpace") + "/" + c.Param("version"),
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
		apiResult, err := executorClient.GetApiResult(ctx, &grpc_executor.ApiRequest{
			StageId:   1,
			ServiceId: 1,
			Page:      1,
			PerPage:   10,
		})

		if err != nil {
			return c.JSON(code.HttpCode(), map[string]interface{}{
				"code": code,
				"msg":  code.Message(),
			})
		} else {
			var datas [10]map[string]interface{}

			for i, v := range apiResult.Data {
				json.Unmarshal([]byte(v), &datas[i])
			}

			return c.JSON(code.HttpCode(), map[string]interface{}{
				"page":         apiResult.Page,
				"perPage":      apiResult.PerPage,
				"totalCount":   apiResult.TotalCount,
				"currentCount": apiResult.CurrentCount,
				"datas":        datas,
			})
		}
	}

	return nil
}
