package v1

import (
	"server/pkg/api/handler"

	"github.com/gin-gonic/gin"
)

type MenuRouter struct{}

func (*MenuRouter) InitRouter(Router gin.IRouter, menuHandler handler.MenuHandlerInterface, middleware ...gin.HandlerFunc) {
	menuRouter := Router.Group("/menus", middleware...)
	{
		menuRouter.GET("", menuHandler.GetMenuList)       // 获取路由菜单列表
		menuRouter.PUT("", menuHandler.UpdateMenu)        // 更新路由菜单
		menuRouter.POST("", menuHandler.AddMenu)          // 添加路由菜单
		menuRouter.DELETE("/:id", menuHandler.DeleteMenu) // 删除路由菜单
	}

	{
		menuRouter.GET("/authorized", menuHandler.GetAsyncMenu) // 获取异步路由菜单
	}
}
