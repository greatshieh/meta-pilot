package initialize

import (
	"server/pkg/global"
	"server/pkg/middleware"
	"strings"

	"github.com/gin-gonic/gin"
)

var appMiddlewares []gin.HandlerFunc

func InitMiddleware(middlewares ...gin.HandlerFunc) {
	appMiddlewares = middlewares
}

func installMiddleware(router gin.IRouter) {
	installSystemMiddleware(router)
	installAppMiddleware(router)
}

func installAppMiddleware(router gin.IRouter) {
	for _, m := range appMiddlewares {
		router.Use(m)
	}
}

// 注册中间件
func installSystemMiddleware(router gin.IRouter) {
	defaultMiddleware := defaultMiddleware()
	for _, m := range global.MPA_CONFIG.System.Middleware {
		m = strings.ToLower(m)
		if mw, ok := defaultMiddleware[m]; ok {
			if m == "jwt" {
				middleware.InitJWT()
			}
			router.Use(mw)
		}
	}

	router.Use(middleware.TimeoutMiddleware())
	global.MPA_LOG.Info("🚀 middleware register success")
}

func defaultMiddleware() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"recovery":  middleware.Recovery(true),  // 异常捕捉中间件
		"jwt":       middleware.JWTAuth(),       // jwt认证中间件
		"accesslog": middleware.AccessLog(),     // 操作记录
		"cors":      middleware.CorsByRules(),   // cors
		"casbin":    middleware.CasbinHandler(), // casbin权限认证中间件
	}
}
