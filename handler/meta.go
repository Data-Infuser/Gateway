package handler

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/promptech1/infuser-gateway/enum"
	"gitlab.com/promptech1/infuser-gateway/handler/middleware"
	grpc_executor "gitlab.com/promptech1/infuser-gateway/infuser-protobuf/gen/proto/executor"
)

type Results struct {
	Meta []*grpc_executor.SchemaMeta `json:"cols" xml:"col"`
}

func (h *Handler) FindMeta(c echo.Context) error {
	ctx := c.Request().Context()
	dataType := c.QueryParam("returnType")

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

	if dataType == "XML" {
		tmp := struct {
			Results
			XMLName struct{} `xml:"cols"`
		}{Results: Results{
			Meta: schemaResult.Meta,
		}}
		return c.XML(enum.Valid.HttpCode(), tmp)
	} else {
		return c.JSON(enum.Valid.HttpCode(), &Results{
			Meta: schemaResult.Meta,
		})
	}
}
