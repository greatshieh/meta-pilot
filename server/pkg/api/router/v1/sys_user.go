package v1

import (
	"server/pkg/api/handler"

	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (s *UserRouter) InitRouter(Router gin.IRouter, userHandler handler.UserHandlerInterface, middleware ...gin.HandlerFunc) {
	userRouter := Router.Group("/users", middleware...)

	// 认证相关路由
	userRouter.POST("/register", userHandler.Register)            // 管理员注册账号
	userRouter.POST("/login", userHandler.Login)                  // 用户登录
	userRouter.PATCH("/password", userHandler.ChangePassword)     // 用户修改密码
	userRouter.POST("/password/reset", userHandler.ResetPassword) // 重置用户密码

	// 用户管理路由
	userRouter.GET("", userHandler.GetUserList)                          // 分页获取用户列表
	userRouter.DELETE("/:id", userHandler.DeleteUser)                    // 删除用户
	userRouter.PUT("/:id", userHandler.SetUserInfo)                      // 设置用户信息
	userRouter.PATCH("/:id/authorities", userHandler.SetUserAuthorities) // 设置用户权限列表

	// 当前用户相关路由
	userRouter.GET("/me", userHandler.GetUserInfo)                  // 获取自身信息
	userRouter.PUT("/me", userHandler.SetSelfInfo)                  // 设置自身信息
	userRouter.PATCH("/me/authority", userHandler.SetUserAuthority) // 登录用户切换权限
}
