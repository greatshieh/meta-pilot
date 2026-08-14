package system

import (
	"context"
	systemReq "server/pkg/api/request"
	systemRes "server/pkg/api/response"
	"server/pkg/errcode"
	"server/pkg/global"
	"server/pkg/model/system"
	"server/pkg/utils"

	"github.com/go-redis/redis/v8"
	"github.com/marmotedu/errors"
)

// AddTokenToBlacklist 将 token 添加到黑名单
// 参数：
//
//	token: 要添加到黑名单的 JWT token
//
// 返回值：
//
//	error: 错误信息
//
// 功能：
//
//	创建 JWT 黑名单记录并保存到数据库，使该 token 失效
func (j *JwtService) AddTokenToBlacklist(ctx context.Context, token string) error {
	jwt := system.JwtBlacklist{Jwt: token}

	return j.jwtDao.JsonInBlacklist(ctx, jwt)
}

// GenerateToken 生成 JWT 令牌
// 参数：
//
//	user: 用户信息
//
// 返回值：
//
//	systemRes.LoginResponse: 登录响应信息，包含用户信息、token 和过期时间
//	error: 错误信息
//
// 功能：
//
//	为用户生成 JWT 令牌，并根据系统配置处理多点登录限制
func (j *JwtService) GenerateToken(ctx context.Context, user system.SysUser) (systemRes.LoginResponse, error) {
	jwt := &utils.JWT{SigningKey: []byte(global.MPA_CONFIG.JWT.SigningKey)} // 唯一签名

	claims := jwt.CreateClaims(systemReq.BaseClaims{
		UUID:        user.UUID,
		ID:          user.ID,
		NickName:    user.NickName,
		UserName:    user.UserName,
		AuthorityId: user.AuthorityID,
	})

	token, err := jwt.CreateToken(claims)

	if err != nil {
		return systemRes.LoginResponse{}, err
	}

	if !global.MPA_CONFIG.System.UseMultipoint {
		// 更新最近一次登录时间
		if err := j.userDao.UpdateLoginTime(ctx, user.ID); err != nil {
			return systemRes.LoginResponse{}, err
		}

		return systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, nil
	}

	// 开启多点在线限制, 从 Redis 中获取用户 Token
	if jwtStr, err := j.jwtDao.GetRedisJWT(ctx, user.UserName); err == redis.Nil {
		// Redis 中没有保存用户 Token, 写入新 Token 到 Redis
		if err := j.jwtDao.SetRedisJWT(ctx, token, user.UserName); err != nil {
			return systemRes.LoginResponse{}, err
		}

		// 更新最近一次登录时间
		if err := j.userDao.UpdateLoginTime(ctx, user.ID); err != nil {
			return systemRes.LoginResponse{}, err
		}

		return systemRes.LoginResponse{
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, nil
	} else if err != nil {
		// 读取 Redis 中的 Token 失败
		return systemRes.LoginResponse{}, errors.WithCode(errcode.ErrTokenRedisStateSetFailed, err.Error())
	} else {
		// 从 Redis 中获取了 Token
		var blackJWT system.JwtBlacklist
		blackJWT.Jwt = jwtStr

		// 将旧 Token 加入黑名单
		if err := j.jwtDao.JsonInBlacklist(ctx, blackJWT); err != nil {
			return systemRes.LoginResponse{}, errors.WithCode(errcode.ErrExpired, "%s", err.Error())
		}

		// 从新 Token 覆盖 Redis 中的旧 Token
		if err := j.jwtDao.SetRedisJWT(ctx, token, user.UserName); err != nil {
			return systemRes.LoginResponse{}, err
		}

		// 更新最近一次登录时间
		if err := j.userDao.UpdateLoginTime(ctx, user.ID); err != nil {
			return systemRes.LoginResponse{}, err
		}

		return systemRes.LoginResponse{
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, nil
	}
}
