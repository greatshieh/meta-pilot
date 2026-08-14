package initialize

import (
	"server/pkg/global"
	"server/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func InitstallRateLimiter(router gin.IRouter) {
	for _, rl := range global.MPA_CONFIG.RateLimiter {
		router.Use(middleware.RateLimiter(rl.Type, rl.Rules...))
	}
}
