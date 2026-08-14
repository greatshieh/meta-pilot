package v1

import (
	"server/pkg/api/handler"

	"github.com/gin-gonic/gin"
)

type CasbinRouter struct{}

func (s *CasbinRouter) InitRouter(Router gin.IRouter, casbinHandler handler.CasbinHandlerInterface, middleware ...gin.HandlerFunc) {
	casbinRouter := Router.Group("/policies", middleware...)

	{
		casbinRouter.POST("/authorities/:authorityId", casbinHandler.SetCasbin)                 // 设置指定角色的权限策略
		casbinRouter.GET("/authorities/:authorityId", casbinHandler.GetPolicyPathByAuthorityId) // 获取指定角色的权限路径
	}
}
