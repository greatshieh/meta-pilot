package v1

import (
	"net/http"
	"server/internal/api"
	"server/pkg/model/common/response"
	sysUtils "server/pkg/utils"

	"github.com/gin-gonic/gin"
)

type WechatHandler struct{}

func (*WechatHandler) Get(c *gin.Context) {
	var query api.ServerValidateReqModel

	if err := sysUtils.VerifyBindQuery(c, &query); err != nil {
		response.FailResponse(c, err)
		return
	}

	echoStr := wechatService.Get(query)

	if echoStr != "" {
		c.Writer.Write([]byte(query.Echostr))
		return
	}
	c.JSON(http.StatusOK, "")
}
