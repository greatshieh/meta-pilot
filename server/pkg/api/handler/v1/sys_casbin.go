package v1

import (
	"fmt"
	"server/pkg/api/handler"
	"server/pkg/api/request"
	"server/pkg/errcode"
	"server/pkg/model/common/response"
	"server/pkg/service"
	"server/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
)

// CasbinHandler Casbin处理器结构体
// 处理权限策略相关的HTTP请求
type CasbinHandler struct {
	casbinService service.CasbinServiceInterface
}

// NewCasbinHandler 创建Casbin处理器实例
func NewCasbinHandler(casbinService service.CasbinServiceInterface) handler.CasbinHandlerInterface {
	return &CasbinHandler{
		casbinService: casbinService,
	}
}

// SetCasbin 设置Casbin权限策略
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理设置Casbin权限策略请求，验证请求参数，调用服务层设置权限策略，返回设置结果
func (cas *CasbinHandler) SetCasbin(c *gin.Context) {
	authorityID := c.Param("authorityId")
	if authorityID == "" {
		response.FailResponse(c, errors.WithCode(errcode.ErrValidation, "authorityId is empty"))
		return
	}
	authorityIDInt, _ := strconv.Atoi(authorityID)
	authorityIDUint := uint(authorityIDInt)

	adminAuthorityID := utils.GetUserAuthorityId(c)

	var cmr request.CasbinInReceive
	if err := utils.VerifyBindJson(c, &cmr); err != nil {
		response.FailResponse(c, err)
		return
	}

	if err := cas.casbinService.SetCasbin(c.Request.Context(), adminAuthorityID, authorityIDUint, cmr.CasbinInfos); err != nil {
		response.FailResponse(c, errors.WithCode(errcode.ErrCasbinSet, "%s", fmt.Sprintf("%d 权限设置失败", authorityIDUint)))

		return
	}
	response.SuccessResponse(c, nil)
}

// GetPolicyPathByAuthorityId 获取权限路径
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理获取权限路径请求，验证请求参数，调用服务层获取权限路径，返回路径数据
func (cas *CasbinHandler) GetPolicyPathByAuthorityId(c *gin.Context) {
	var authorityID = c.Param("authorityId")

	paths := cas.casbinService.GetPolicyPathByAuthorityID(c.Request.Context(), authorityID)
	response.SuccessResponse(c, paths)
}
