package v1

import (
	"server/pkg/api/handler"
	"server/pkg/api/request"
	"server/pkg/global"
	"server/pkg/model/common/response"
	"server/pkg/model/system"
	"server/pkg/service"
	"server/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AuthorityHandler 权限处理器结构体
// 处理权限相关的 HTTP 请求
type AuthorityHandler struct {
	authorityService service.AuthorityServiceInterface
}

// NewAuthorityHandler 创建权限处理器实例
func NewAuthorityHandler(authorityService service.AuthorityServiceInterface) handler.AuthorityHandlerInterface {
	return &AuthorityHandler{
		authorityService: authorityService,
	}
}

// GetAuthorityList 获取权限列表
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理获取权限列表请求，根据用户权限返回相应的权限数据
func (a *AuthorityHandler) GetAuthorityList(c *gin.Context) {
	// 从请求上下文中获取当前用户的权限ID
	authorityID := utils.GetUserAuthorityId(c)

	// 调用权限服务获取权限列表，传入当前用户的权限ID作为参数
	// 这样可以确保用户只能看到自己有权限访问的数据
	list, err := a.authorityService.GetAuthorityList(c.Request.Context(), authorityID)

	// 如果获取权限列表过程中发生错误，则返回错误响应
	if err != nil {
		// 使用自定义错误码和错误信息封装错误，并返回失败响应
		response.FailResponse(c, err)
		return
	}

	// 获取权限列表成功，返回成功响应并将权限列表数据包含在响应中
	response.SuccessResponse(c, list)
}

// AddAuthority 添加权限
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理添加权限请求，验证请求参数，调用服务层添加权限，返回添加结果
func (a *AuthorityHandler) AddAuthority(c *gin.Context) {
	// 定义接收请求参数的权限结构体
	var authority system.SysRoleAuthority

	// 绑定JSON格式的请求参数到权限结构体
	if err := utils.VerifyBindJson(c, &authority, utils.AuthorityVerify); err != nil {
		// 参数绑定失败，返回带错误码的错误响应
		response.FailResponse(c, err)
		return
	}

	// 在严格权限模式下，如果要创建的是顶级权限（父ID为0），则将其父ID设置为当前用户的权限ID
	// 这样可以确保权限的层级关系正确，防止普通用户创建与自己无关的顶级权限
	if *authority.ParentID == 0 && global.MPA_CONFIG.System.UseStrictAuth {
		// 将当前用户的权限ID作为新建权限的父ID
		authority.ParentID = utils.Pointer(utils.GetUserAuthorityId(c))
	}

	// 调用服务层方法添加权限
	auth, err := a.authorityService.AddAuthority(c.Request.Context(), authority)
	if err != nil {
		// 添加失败，返回错误响应
		response.FailResponse(c, err)
		return
	}

	// 添加成功，返回成功的响应和创建的权限信息
	response.SuccessResponse(c, auth)
}

// DeleteAuthority 删除权限
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理删除权限请求，验证请求参数，调用服务层删除权限，返回删除结果
func (a *AuthorityHandler) DeleteAuthority(c *gin.Context) {
	// 定义接收请求参数的权限结构体
	var authority system.SysRoleAuthority

	// 绑定JSON格式的请求参数到权限结构体
	if err := utils.VerifyBindJson(c, &authority, utils.AuthorityVerify); err != nil {
		// 参数绑定失败，返回带错误码的错误响应
		response.FailResponse(c, err)
		return
	}

	// 删除角色之前需要判断是否有用户正在使用此角色
	if err := a.authorityService.DeleteAuthority(c.Request.Context(), authority); err != nil {
		response.FailResponse(c, err)
		return
	}
	// _ = casbinService.FreshCasbin()
	response.SuccessResponse(c, "OK")
}

// UpdateAuthority 更新权限
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理更新权限请求，验证请求参数，调用服务层更新权限，返回更新结果
func (a *AuthorityHandler) UpdateAuthority(c *gin.Context) {
	var auth system.SysRoleAuthority

	if err := utils.VerifyBindJson(c, &auth, utils.AuthorityVerify); err != nil {
		response.FailResponse(c, err)
		return
	}

	authority, err := a.authorityService.UpdateAuthority(c.Request.Context(), auth)
	if err != nil {
		response.FailResponse(c, err)
		return
	}

	response.SuccessResponse(c, authority)
}

// GetAuthorityMenu 获取权限菜单
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理获取权限菜单请求，验证请求参数，调用服务层获取权限菜单，返回菜单数据
func (a *AuthorityHandler) GetAuthorityMenu(c *gin.Context) {
	authorityID := c.Param("authorityId")

	authorityIDInt, _ := strconv.Atoi(authorityID)

	// 调用权限服务的GetAuthorityMenu方法获取指定权限ID对应的菜单列表
	menus, err := a.authorityService.GetAuthorityMenu(c.Request.Context(), uint(authorityIDInt))

	// 检查获取菜单权限列表过程中是否出现错误
	if err != nil {
		// 如果获取过程中出现错误，使用统一响应格式返回错误信息
		response.FailResponse(c, err)
		return
	}

	// 如果操作成功，返回成功的响应信息
	response.SuccessResponse(c, menus)

}

// SetMenuAuthority 设置菜单权限
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理设置菜单权限请求，验证请求参数，调用服务层设置菜单权限，返回设置结果
func (a *AuthorityHandler) SetMenuAuthority(c *gin.Context) {
	authorityID := c.Param("authorityId")

	authorityIDInt, _ := strconv.Atoi(authorityID)

	// 定义接收请求参数的结构体
	var authorityMenu request.SetMenuAuthority

	// 绑定JSON格式的请求参数到结构体
	if err := utils.VerifyBindJson(c, &authorityMenu); err != nil {
		// 参数绑定失败，返回错误响应
		response.FailResponse(c, err)
		return
	}

	// 从上下文中获取当前用户的权限ID（操作者权限）
	adminAuthorityID := utils.GetUserAuthorityId(c)

	// 调用权限服务方法设置菜单权限
	if err := a.authorityService.SetMenuAuthority(c.Request.Context(), authorityMenu, uint(authorityIDInt), adminAuthorityID); err != nil {
		response.FailResponse(c, err)
		return
	}

	// 设置成功，返回成功响应
	response.SuccessResponse(c, nil, "权限设置成功")
}
