package v1

import (
	"server/pkg/api/handler"

	"github.com/gin-gonic/gin"
)

type ApiRouter struct{}

func NewApiRouter() ApiRouter {
	return ApiRouter{}
}

func (*ApiRouter) InitRouter(Router gin.IRouter, apiHandler handler.ApiHandlerInterface, middleware ...gin.HandlerFunc) {
	apiRouter := Router.Group("/apis", middleware...)
	{
		apiRouter.GET("", apiHandler.GetApiList)      // 分页获取所有 api 列表
		apiRouter.POST("", apiHandler.CreateApi)      // 创建新的 api
		apiRouter.DELETE(":id", apiHandler.DeleteApi) // 删除 api
		apiRouter.PUT("", apiHandler.UpdateApi)       // 修改 api 属性
	}

	{
		apiRouter.GET("/all", apiHandler.GetAllApis)      // 获取所有 api 列表
		apiRouter.GET("/groups", apiHandler.GetApiGroups) // 获取路由组
		apiRouter.POST("/ignore", apiHandler.IgnoreApi)   // 忽略Api
		apiRouter.GET("/sync", apiHandler.SyncApi)        // 获取待同步Api
		apiRouter.PUT("/sync", apiHandler.EnterSyncApi)   // 确认同步Api
	}
}
