package v1

import (
	v1 "server/internal/api/handler/v1"

	"github.com/gin-gonic/gin"
)

type WeChatRouter struct{}

func (w *WeChatRouter) InitRouter(Router gin.IRouter, middleware ...gin.HandlerFunc) {
	group := Router.Group("wechat", middleware...)

	wechatHandler := new(v1.WechatHandler)
	group.GET("", wechatHandler.Get)
}
