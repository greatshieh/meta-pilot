package middleware

import (
	"fmt"
	"server/pkg/errcode"
	"server/pkg/global"
	"server/pkg/model/common/response"

	"server/pkg/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
)

// CasbinHandler 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取请求的PATH
		path := c.Request.URL.Path
		obj := strings.TrimPrefix(path, global.MPA_CONFIG.System.RouterPrefix)
		// 获取请求方法
		act := c.Request.Method
		// 检查请求的路径是否在被忽略的API列表中
		for _, ignoredAPI := range global.MPA_IGNORED_APIS {
			if ignoredAPI.Path == obj && ignoredAPI.Method == act {
				c.Next()
			}
		}

		waitUse, _ := utils.GetClaims(c)

		// 获取用户的角色
		sub := strconv.Itoa(int(waitUse.AuthorityId))
		e := utils.GetCasbin() // 判断策略中是否存在
		success, _ := e.Enforce(sub, obj, act)
		if !success {
			response.FailResponse(c, errors.WithCode(errcode.ErrPermissionDenied, fmt.Sprintf("casbin 权限不足, path: %s, act: %s, sub: %s", obj, act, sub)))
			c.Abort()
			return
		}
		c.Next()
	}
}
