package v1

import (
	"server/pkg/api/handler"

	"github.com/gin-gonic/gin"
)

type AuthorityRouter struct{}

func (*AuthorityRouter) InitRouter(Router gin.IRouter, authorityHandler handler.AuthorityHandlerInterface, middleware ...gin.HandlerFunc) {
	authorityRouter := Router.Group("/authorities", middleware...)

	{
		authorityRouter.GET("", authorityHandler.GetAuthorityList)   // 获取所有角色列表
		authorityRouter.POST("", authorityHandler.AddAuthority)      // 创建新的角色
		authorityRouter.DELETE("", authorityHandler.DeleteAuthority) // 删除角色
		authorityRouter.PUT("", authorityHandler.UpdateAuthority)    // 修改角色属性
	}

	{
		authorityRouter.GET("/:authorityId/menus", authorityHandler.GetAuthorityMenu)  // 获取指定角色对应的权限菜单
		authorityRouter.POST("/:authorityId/menus", authorityHandler.SetMenuAuthority) // 添加角色菜单权限
	}
}
