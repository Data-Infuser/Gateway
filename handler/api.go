// handler: Rest API 요청 처리 및 gRPC 통신을 통한 데이터 교환 수행 package
package handler

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/clbanning/mxj/v2"
	"github.com/labstack/echo/v4"
	"gitlab.com/promptech1/infuser-gateway/constant"
	"gitlab.com/promptech1/infuser-gateway/enum"
	grpc_author "gitlab.com/promptech1/infuser-gateway/infuser-protobuf/gen/proto/author"
	grpc_executor "gitlab.com/promptech1/infuser-gateway/infuser-protobuf/gen/proto/executor"
)

// ExecuteAPI: 활용자의 Rest API 요청을 처리함. gRPC를 통해 필요한 데이터를 교환하고 그결과를 JSON 형태로 반환함
func (h *Handler) ExecuteAPI(c echo.Context) error {
	ctx := c.Request().Context()
	dataType := c.QueryParam("returnType")

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
		if dataType == "" || dataType == "JSON" {
			return c.JSON(code.HttpCode(), map[string]interface{}{
				"code": code,
				"msg":  code.Message(),
			})
		} else {
			return c.XML(code.HttpCode(), map[string]interface{}{
				"code": code,
				"msg":  code.Message(),
			})
		}
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
	}

	var page int32
	var perPage int32
	cond := map[string]string{}

	// 공공데이터 포털 테스트 인증키 처리 추가
	if token != constant.DATA_PORTAL_TEST_KEY {
		pageInt, err := strconv.ParseInt(c.QueryParam("page"), 10, 32)
		if err != nil {
			page = 1
		} else {
			page = int32(pageInt)
		}

		perPageInt, err := strconv.ParseInt(c.QueryParam("perPage"), 10, 32)
		if err != nil {
			perPage = 10
		} else {
			perPage = int32(perPageInt)
		}

		for key, val := range c.QueryParams() {
			if key[0:4] == "cond" {
				cond[key[5:len(key)-1]] = val[0]
			}
		}
	} else {
		page = 1
		perPage = 10
	}

	if len(apiAuthRes.ProxyEndpoint) > 0 {
		h.Rp(c, apiAuthRes.ProxyEndpoint)
		return nil
	}

	apiResult, err := executorClient.GetApiResult(ctx, &grpc_executor.ApiRequest{
		StageId:   int32(apiAuthRes.AppId),
		ServiceId: int32(apiAuthRes.OperationId),
		Page:      page,
		PerPage:   perPage,
		Cond:      cond,
	})

	if err != nil {
		return c.JSON(code.HttpCode(), map[string]interface{}{
			"code": code,
			"msg":  code.Message(),
		})
	} else {
		if dataType == "XML" {
			var result = map[string]interface{}{
				"page":         apiResult.Page,
				"perPage":      apiResult.PerPage,
				"totalCount":   apiResult.TotalCount,
				"currentCount": apiResult.CurrentCount,
				"matchCount":   apiResult.MatchCount,
			}

			var tmp map[string]interface{}
			var xmlValues []string

			for _, v := range apiResult.Data {
				json.Unmarshal([]byte(v), &tmp)
				for k := range tmp {
					if tmp[k] == nil {
						tmp[k] = ""
					}
				}
				mv := mxj.Map(tmp)
				xmlValue, _ := mv.Xml("item")
				xmlValues = append(xmlValues, string(xmlValue))
			}

			result["data"] = strings.Join(xmlValues[:], "")

			mv := mxj.Map(result)

			xmlResults, err := mv.Xml("results")
			if err != nil {
				h.ctx.Logger.WithField(
					"return err", err.Error(),
				).Error("XML Marshal result")
			}

			c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationXMLCharsetUTF8)
			return c.String(code.HttpCode(), string(xmlResults))
		} else {
			data := make([]map[string]interface{}, apiResult.CurrentCount)

			for i, v := range apiResult.Data {
				json.Unmarshal([]byte(v), &data[i])
			}

			return c.JSON(code.HttpCode(), map[string]interface{}{
				"page":         apiResult.Page,
				"perPage":      apiResult.PerPage,
				"totalCount":   apiResult.TotalCount,
				"currentCount": apiResult.CurrentCount,
				"matchCount":   apiResult.MatchCount,
				"data":         data,
			})
		}
	}

	return nil
}
