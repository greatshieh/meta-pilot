package v1

import (
	"server/pkg/api/handler"
	"server/pkg/model/common/response"

	"github.com/gin-gonic/gin"

	"server/pkg/service"
)

// JwtHandler JWT处理器结构体
// 处理JWT相关的HTTP请求
type JwtHandler struct {
	jwtService service.JwtServiceInterface
}

// NewJwtHandler 创建JWT处理器实例
func NewJwtHandler(jwtService service.JwtServiceInterface) handler.JwtHandlerInterface {
	return &JwtHandler{
		jwtService: jwtService,
	}
}

// JsonInBlacklist 将JWT加入黑名单
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理将JWT加入黑名单请求，从请求头获取token，调用服务层将token加入黑名单，返回操作结果
func (j *JwtHandler) JsonInBlacklist(c *gin.Context) {
	token := c.Request.Header.Get("x-token")

	if err := j.jwtService.AddTokenToBlacklist(c.Request.Context(), token); err != nil {
		response.FailResponse(c, err)
		return
	}

	response.SuccessResponse(c, nil, "jwt作废成功")
}
