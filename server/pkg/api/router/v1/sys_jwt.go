package v1

import (
	"server/pkg/api/handler"

	"github.com/gin-gonic/gin"
)

type JwtRouter struct{}

func (s *JwtRouter) InitRouter(Router gin.IRouter, jwtHandler handler.JwtHandlerInterface, middleware ...gin.HandlerFunc) {
	jwtRouter := Router.Group("/jwt", middleware...)

	jwtRouter.POST("/blacklist", jwtHandler.JsonInBlacklist) // jwt加入黑名单

}
