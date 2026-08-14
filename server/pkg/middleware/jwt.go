package middleware

import (
	v1 "server/pkg/dao/v1"
	"server/pkg/errcode"
	"server/pkg/global"
	"server/pkg/model/common/response"
	"server/pkg/model/system"
	"server/pkg/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/marmotedu/errors"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"go.uber.org/zap"
)

var jwtDao = v1.NewJwtDao(global.MPA_DB)

func InitJWT() {
	dr, err := utils.ParseDuration(global.MPA_CONFIG.JWT.ExpiresTime)
	if err != nil {
		panic(err)
	}
	_, err = utils.ParseDuration(global.MPA_CONFIG.JWT.BufferTime)
	if err != nil {
		panic(err)
	}

	global.BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(dr),
	)
	// file, err := os.Open("go.mod")
	// if err == nil && global.MPA_CONFIG.AutoCode.Module == "" {
	// 	scanner := bufio.NewScanner(file)
	// 	scanner.Scan()
	// 	global.MPA_CONFIG.AutoCode.Module = strings.TrimPrefix(scanner.Text(), "module ")
	// }
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.URL.Path, "/api/sysuser/login") {
			// 如果是登录请求，直接放行
			c.Next()
			return
		}
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.Request.Header.Get("x-token")
		if token == "" {
			response.FailResponse(c, errors.WithCode(errcode.ErrPermissionDenied, "未登录或非法访问"))
			c.Abort()
			return
		}

		if jwtDao.IsBlacklist(c.Request.Context(), token) {
			// Token在黑名单中, 说明这次请求的 Token 已经被新的 Token 替代
			response.FailResponse(c, errors.WithCode(errcode.ErrExpired, "帐户异地登陆或令牌失效"))
			c.Abort()
			return
		}

		j := utils.NewJWT()

		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, utils.ErrTokenExpired) {
				response.FailResponse(c, errors.WithCode(errcode.ErrExpired, "授权已过期"))
				c.Abort()
				return
			}
			response.FailResponse(c, errors.WithCode(errcode.ErrTokenParseFailed, err.Error()))
			c.Abort()
			return
		}

		// 已登录用户被管理员禁用 需要使该用户的jwt失效 此处比较消耗性能 如果需要 请自行打开
		// 用户被删除的逻辑 需要优化 此处比较消耗性能 如果需要 请自行打开

		//if user, err := userDao.FindUserByUuid(claims.UUID.String()); err != nil || user.Enable == 2 {
		//	_ = jwtDao.JsonInBlacklist(system.JwtBlacklist{Jwt: token})
		//	response.FailWithDetailed(gin.H{"reload": true}, err.Error(), c)
		//	c.Abort()
		//}

		if claims.ExpiresAt.Unix()-time.Now().Unix() < claims.BufferTime {
			// Token没有过期, 并且还在缓冲期内, 需要刷新 Token
			dr, _ := utils.ParseDuration(global.MPA_CONFIG.JWT.ExpiresTime)
			// 更新JWT Claims 的过期时间
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(dr))
			// 刷新 Token
			newToken, _ := j.CreateTokenByOldToken(token, *claims)

			newClaims, _ := j.ParseToken(newToken)

			// 更新请求头的 Token
			c.Header("x-token", newToken)
			// 更新 Token 过期时间
			c.Header("x-expires-at", strconv.FormatInt(newClaims.ExpiresAt.Unix(), 10))

			if global.MPA_CONFIG.System.UseMultipoint {
				// 开启了多点在线限制, 将 Redis 中保存的 Token 放入黑名单, 刷新后的 Token 写入 Redis
				RedisJwtToken, err := jwtDao.GetRedisJWT(c.Request.Context(), newClaims.UserName)

				if err != nil {
					// 从 Redis 中取 Token 发生错误
					global.MPA_LOG.Error("get redis jwt failed", zap.Error(err))
				} else {
					// 从 Redis 中取到 Token 后, 将原来的 Token 放进黑名单
					_ = jwtDao.JsonInBlacklist(c.Request.Context(), system.JwtBlacklist{Jwt: RedisJwtToken})
				}

				// 用新的 Token 替换旧 Token
				_ = jwtDao.SetRedisJWT(c.Request.Context(), newToken, newClaims.UserName)
			}
		}
		c.Set("claims", claims)
		c.Next()
	}
}
