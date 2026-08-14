package initialize

import (
	"net/http"
	"server/pkg/api/router"
	"server/pkg/global"

	"github.com/gin-gonic/gin"
)

var appRouter []router.Router

func installRouter(router *gin.RouterGroup, customRouters ...router.Router) {
	// 注册自定义路由
	for _, r := range customRouters {
		r.InitRouter(router)
	}
	// 安装系统路由
	installtSystemRouter(router)
	// 安装应用路由
	installAppRouter(router)

	global.MPA_LOG.Info("🚀 路由表注册成功")
}

func installAppRouter(router gin.IRouter) {
	for _, r := range appRouter {
		r.InitRouter(router)
	}
}

// 初始化系统路由
func installtSystemRouter(r gin.IRouter) {
	// 健康监测
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	// 添加系统路由前缀
	routerEenter := r.Group(global.MPA_CONFIG.System.SystemRouterPrefix)

	var sysRouter = router.SysRouter{}
	sysRouter.InitRouter(routerEenter)
}
