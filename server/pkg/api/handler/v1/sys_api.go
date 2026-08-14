package v1

import (
	"server/pkg/api/handler"
	"server/pkg/api/request"
	"server/pkg/model/common/response"
	"server/pkg/model/system"
	"server/pkg/service"
	"server/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"

	sysUtils "server/pkg/utils"
)

// ApiHandler API 处理器结构体
// 处理 API 相关的 HTTP 请求
type ApiHandler struct {
	apiService service.ApiServiceInterface
}

// NewApiHandler 创建 API 处理器实例
func NewApiHandler(apiService service.ApiServiceInterface) handler.ApiHandlerInterface {
	return &ApiHandler{
		apiService: apiService,
	}
}

// CreateApi 创建 API
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理创建API请求，验证请求参数，调用服务层创建API，返回创建结果
func (a *ApiHandler) CreateApi(c *gin.Context) {
	var api system.SysApi

	if err := sysUtils.VerifyBindJson(c, &api, sysUtils.ApiVerify); err != nil {
		response.FailResponse(c, err)
		return
	}

	api, err := a.apiService.CreateAPI(c.Request.Context(), api)
	if err != nil {
		response.FailResponse(c, err)
		return
	}

	response.SuccessResponse(c, api)
}

// GetApiList 获取 API 列表
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理获取API列表请求，验证请求参数，调用服务层获取API列表，返回分页结果
func (a *ApiHandler) GetApiList(c *gin.Context) {
	var pageInfo request.SearchApiParams
	if err := sysUtils.VerifyBindQuery(c, &pageInfo, sysUtils.PageInfoVerify); err != nil {
		response.FailResponse(c, err)
		return
	}

	result, err := a.apiService.GetAPIList(c.Request.Context(), pageInfo)
	if err != nil {
		response.FailResponse(c, err)
		return
	}
	response.SuccessResponse(c, result)
}

// GetAllApis 获取所有 API
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理获取所有API请求，调用服务层获取所有API，返回API列表
func (a *ApiHandler) GetAllApis(c *gin.Context) {
	authorityID := utils.GetUserAuthorityId(c)
	apis, err := a.apiService.GetAllAPIs(c.Request.Context(), authorityID)
	if err != nil {
		response.FailResponse(c, err)
		return
	}
	response.SuccessResponse(c, apis)
}

// DeleteApi 删除 API
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理删除API请求，验证请求参数，调用服务层删除API，返回删除结果
func (a *ApiHandler) DeleteApi(c *gin.Context) {
	reqID := c.Param("id")

	id, _ := strconv.Atoi(reqID)

	if err := a.apiService.DeleteAPI(c.Request.Context(), uint(id)); err != nil {
		response.FailResponse(c, err)
		return
	}

	response.SuccessResponse(c, nil, "删除成功")
}

// UpdateApi 更新 API
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理更新API请求，验证请求参数，调用服务层更新API，返回更新结果
func (a *ApiHandler) UpdateApi(c *gin.Context) {
	var api system.SysApi

	if err := sysUtils.VerifyBindJson(c, &api, sysUtils.ApiVerify); err != nil {
		response.FailResponse(c, err)
		return
	}

	api, err := a.apiService.UpdateAPI(c.Request.Context(), api)
	if err != nil {
		response.FailResponse(c, err)
		return
	}

	response.SuccessResponse(c, api, "API 更新成功")
}

// GetApiGroups 获取 API 分组
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理获取API分组请求，调用服务层获取API分组，返回分组列表
func (a *ApiHandler) GetApiGroups(c *gin.Context) {
	groups, err := a.apiService.GetAPIGroups(c.Request.Context())
	if err != nil {
		response.FailResponse(c, err)
		return
	}
	response.SuccessResponse(c, groups)
}

// SyncApi 同步 API
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理同步API请求，调用服务层同步API，返回同步结果
func (a *ApiHandler) SyncApi(c *gin.Context) {
	newApis, deleteApis, ignoreApis, err := a.apiService.SyncAPI(c.Request.Context())
	if err != nil {
		response.FailResponse(c, err)
		return
	}
	response.SuccessResponse(c, gin.H{
		"newApis":    newApis,
		"deleteApis": deleteApis,
		"ignoreApis": ignoreApis,
	})
}

// IgnoreApi 忽略 API
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理忽略API请求，验证请求参数，调用服务层忽略API，返回忽略结果
func (a *ApiHandler) IgnoreApi(c *gin.Context) {
	var ignoreApi system.SysIgnoreApi
	if err := sysUtils.VerifyBindJson(c, &ignoreApi); err != nil {
		response.FailResponse(c, err)
		return
	}

	if err := a.apiService.IgnoreAPI(c.Request.Context(), ignoreApi); err != nil {
		response.FailResponse(c, err)
		return
	}
	response.SuccessResponse(c, nil)
}

// EnterSyncApi 进入同步 API
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理进入同步API请求，验证请求参数，调用服务层进入同步API，返回进入结果
func (a *ApiHandler) EnterSyncApi(c *gin.Context) {
	var enterSync request.EnterSyncApiParams
	if err := sysUtils.VerifyBindJson(c, &enterSync); err != nil {
		response.FailResponse(c, err)
		return
	}

	if err := a.apiService.EnterSyncAPI(c.Request.Context(), enterSync); err != nil {
		response.FailResponse(c, err)
		return
	}
	response.SuccessResponse(c, nil)
}
