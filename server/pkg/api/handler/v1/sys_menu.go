package v1

import (
	"server/pkg/api/handler"
	"server/pkg/model/common/response"
	"server/pkg/model/system"
	"server/pkg/service"
	"server/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MenuHandler 菜单处理器结构体
// 处理菜单相关的HTTP请求
type MenuHandler struct {
	menuService service.MenuServiceInterface
}

// NewMenuHandler 创建菜单处理器实例
func NewMenuHandler(menuService service.MenuServiceInterface) handler.MenuHandlerInterface {
	return &MenuHandler{
		menuService: menuService,
	}
}

// GetAsyncMenu 获取异步菜单
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理获取异步菜单请求，根据用户权限返回相应的菜单数据
func (m *MenuHandler) GetAsyncMenu(c *gin.Context) {
	// 从请求上下文中获取当前用户的活动权限ID
	authorityID := utils.GetUserAuthorityId(c)

	// 调用菜单服务获取异步菜单数据，传入当前用户的权限ID作为参数
	// 这样可以确保用户只能看到自己有权限访问的菜单项
	menus, err := m.menuService.GetAsyncMenu(c.Request.Context(), authorityID)

	// 如果获取菜单数据过程中发生错误，则返回错误响应
	if err != nil {
		// 将错误信息写入响应并返回，状态标记为"failure"
		response.FailResponse(c, err)
		return
	}

	// 获取菜单数据成功，返回成功响应并将菜单数据包含在响应中
	response.SuccessResponse(c, menus)
}

// GetMenuList 获取菜单列表
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理获取菜单列表请求，根据用户权限返回完整的菜单树结构
func (m *MenuHandler) GetMenuList(c *gin.Context) {
	// 从请求上下文获取当前用户的权限ID
	// GetUserAuthorityId会解析请求中的JWT token来获取用户权限信息
	authorityID := utils.GetUserAuthorityId(c)

	// 调用菜单服务的GetMenuList方法，根据权限ID获取菜单列表
	// GetMenuList会构建完整的菜单树结构，包含所有层级的子菜单
	menus, err := m.menuService.GetMenuList(c.Request.Context(), authorityID)

	// 检查获取菜单列表过程中是否出现错误
	if err != nil {
		response.FailResponse(c, err)
		return
	}

	// 使用统一响应格式返回成功的菜单列表数据
	response.SuccessResponse(c, menus)
}

// UpdateMenu 更新菜单信息
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理更新菜单请求，验证请求参数，调用服务层更新菜单，返回更新结果
func (m *MenuHandler) UpdateMenu(c *gin.Context) {
	// 声明一个SysMenu类型的变量，用于接收解析后的菜单数据
	var menu system.SysMenu

	// 使用ShouldBindJSON方法将请求体中的JSON数据绑定到menu变量
	if err := utils.VerifyBindJson(c, &menu); err != nil {
		response.FailResponse(c, err)
		return
	}

	// 调用菜单服务的UpdateMenu方法执行菜单更新操作
	// 将解析后的菜单对象传递给服务层进行数据库更新
	if err := m.menuService.UpdateMenu(c.Request.Context(), menu); err != nil {
		// 如果更新过程中出现错误，使用统一响应格式返回错误信息
		response.FailResponse(c, err)
		return
	}

	// 如果更新成功，返回成功的响应信息
	response.SuccessResponse(c, nil, "菜单更新成功")
}

// AddMenu 添加新菜单
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理添加新菜单请求，验证请求参数，调用服务层创建菜单，返回创建结果
func (m *MenuHandler) AddMenu(c *gin.Context) {
	// 声明一个SysMenu类型的变量，用于接收解析后的菜单数据
	var menu system.SysMenu

	if err := utils.VerifyBindJson(c, &menu, utils.MenuVerify, utils.MenuMetaVerify); err != nil {
		// 如果JSON解析失败，使用统一响应格式返回错误信息
		response.FailResponse(c, err)
		return
	}

	// 调用菜单服务的AddMenu方法执行菜单创建操作
	// 将解析后的菜单对象传递给服务层进行数据库插入
	if err := m.menuService.AddMenu(c.Request.Context(), menu); err != nil {
		// 如果创建过程中出现错误，使用统一响应格式返回错误信息
		response.FailResponse(c, err)
		return
	}

	// 如果创建成功，返回成功的响应信息
	response.SuccessResponse(c, nil, "菜单创建成功")
}

// DeleteMenu 删除菜单
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理删除菜单请求，验证请求参数，调用服务层删除菜单，返回删除结果
func (m *MenuHandler) DeleteMenu(c *gin.Context) {
	// 从请求URL中提取菜单ID参数
	id := c.Param("id")

	menuID, _ := strconv.ParseUint(id, 10, 64)

	// 调用菜单服务的DeleteMenu方法执行菜单删除操作
	// 将提取的菜单ID作为参数传递给服务层
	if err := m.menuService.DeleteMenu(c.Request.Context(), uint(menuID)); err != nil {
		// 如果删除过程中出现错误，使用统一响应格式返回错误信息
		response.FailResponse(c, err)
		return
	}

	// 如果删除成功，返回成功的响应信息
	response.SuccessResponse(c, nil, "菜单删除成功")
}
