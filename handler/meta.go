package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"gitlab.com/promptech1/infuser-gateway/enum"
	"gitlab.com/promptech1/infuser-gateway/handler/middleware"
	grpc_executor "gitlab.com/promptech1/infuser-gateway/infuser-protobuf/gen/proto/executor"
)

func (h *Handler) FindMeta(c echo.Context) error {
	ctx := c.Request().Context()

	authRes, authCode := middleware.CheckAuth(h.authPool, c)
	if !authCode.Valid() {
		return c.JSON(authCode.HttpCode(), map[string]interface{}{
			"code": authCode,
			"msg":  authCode.Message(),
		})
	}

	executorConn, err := h.executorPool.Get(ctx)
	defer executorConn.Close()
	if err != nil {
		code := enum.InternalException
		return c.JSON(code.HttpCode(), map[string]interface{}{
			"code": code,
			"msg":  code.Message(),
		})
	}

	executorClient := grpc_executor.NewSchemaMetaServiceClient(executorConn)
	schemaResult, err := executorClient.FindSchemaMeta(ctx, &grpc_executor.SchemaMetaReq{
		StageId:   int32(authRes.AppId),
		ServiceId: int32(authRes.OperationId),
	})

	fmt.Printf("%+v", schemaResult)

	return c.JSON(enum.Valid.HttpCode(), map[string]interface{}{
		"cols": schemaResult.Meta,
	})
}
