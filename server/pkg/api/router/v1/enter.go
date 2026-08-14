package v1

import (
	"github.com/gin-gonic/gin"
)

type ApiGroup struct {
	ApiRouter
	AuthorityRouter
	CasbinRouter
	JwtRouter
	MenuRouter
	UserRouter
}

func NewApiGroup() ApiGroup {
	return ApiGroup{}
}

// InitRouter 初始化 v1 版本的 API 路由
// 参数：
//
//	r: Gin 路由器实例
//	middleware: 全局中间件列表
//
// 功能：
//
//	创建 v1 路由组
//	初始化各个模块的路由，如微信相关路由
//	为所有 v1 路由应用全局中间件
func InitRouter(r gin.IRouter, middleware ...gin.HandlerFunc) {
	// 创建 v1 版本的路由组
	group := r.Group("v1")

	var api = NewApiGroup()

	// 初始化微信相关路由
	var (
		apiHandler       = newApiHandler()
		authorityHandler = newAuthorityHandler()
		casbinHandler    = newCasbinHandler()
		jwtHandler       = newJwtHandler()
		menuHandler      = newMenuHandler()
		userHandler      = newUserHandler()
	)

	api.ApiRouter.InitRouter(group, apiHandler)
	api.AuthorityRouter.InitRouter(group, authorityHandler)
	api.CasbinRouter.InitRouter(group, casbinHandler)
	api.JwtRouter.InitRouter(group, jwtHandler)
	api.MenuRouter.InitRouter(group, menuHandler)
	api.UserRouter.InitRouter(group, userHandler)

	// 后续可以添加其他模块的路由初始化
	// 例如：s
	// userRouter := new(UserRouter)
	// userRouter.InitRouter(group)
}
