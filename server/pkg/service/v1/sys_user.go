package system

import (
	"context"
	"server/pkg/errcode"
	"server/pkg/global"
	"server/pkg/model/common/response"
	"server/pkg/model/system"
	"server/pkg/utils"

	systemReq "server/pkg/api/request"
	systemRes "server/pkg/api/response"

	"time"

	"github.com/go-redis/redis/v8"
	"github.com/marmotedu/errors"
	uuid "github.com/satori/go.uuid"
)

// Register 用户注册
// 参数：
//
//	req: 注册请求信息
//
// 返回值：
//
//	systemRes.SysUserResponse: 用户响应信息
//	error: 错误信息
func (u *UserService) Register(ctx context.Context, req systemReq.Register) (systemRes.SysUserResponse, error) {
	var authorities []system.SysRoleAuthority
	for _, v := range req.AuthorityIds {
		authorities = append(authorities, system.SysRoleAuthority{
			AuthorityID: v,
		})
	}

	user := &system.SysUser{
		UserName:    req.UserName,
		NickName:    req.NickName,
		AuthorityID: req.AuthorityId,
		Authorities: authorities,
		Password:    req.Password,
		Phone:       req.Phone,
		Email:       req.Email,
	}

	userReturn, err := u.userDao.Register(ctx, *user)
	if err != nil {
		return systemRes.SysUserResponse{}, err
	}

	return systemRes.SysUserResponse{User: userReturn}, nil
}

// Login 用户登录
// 参数：
//
//	l: 登录请求信息
//
// 返回值：
//
//	systemRes.LoginResponse: 登录响应信息
//	error: 错误信息
func (u *UserService) Login(ctx context.Context, l systemReq.Login) (systemRes.LoginResponse, error) {
	user := &system.SysUser{UserName: l.UserName, Password: l.Password}
	userInter, err := u.userDao.Login(ctx, *user)

	if err != nil {
		return systemRes.LoginResponse{}, err
	}

	if !userInter.IsActive {
		return systemRes.LoginResponse{}, errors.WithCode(errcode.ErrUserForbidden, "用户被禁止登录")
	}

	return u.jwtService.GenerateToken(ctx, *userInter)
}

// GenerateToken 生成用户登录令牌
// 参数：
//
//	user: 用户信息
//
// 返回值：
//
//	systemRes.LoginResponse: 登录响应信息，包含用户信息、token和过期时间
//	error: 错误信息
func (u *UserService) GenerateToken(ctx context.Context, user system.SysUser) (systemRes.LoginResponse, error) {
	j := &utils.JWT{SigningKey: []byte(global.MPA_CONFIG.JWT.SigningKey)}

	claims := j.CreateClaims(systemReq.BaseClaims{
		UUID:        user.UUID,
		ID:          user.ID,
		NickName:    user.NickName,
		UserName:    user.UserName,
		AuthorityId: user.AuthorityID,
	})

	token, err := j.CreateToken(claims)
	if err != nil {
		return systemRes.LoginResponse{}, err
	}

	if !global.MPA_CONFIG.System.UseMultipoint {
		global.MPA_DB.Model(&user).UpdateColumn("last_login", time.Now().Local())

		return systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, nil
	}

	if jwtStr, err := u.jwtDao.GetRedisJWT(ctx, user.UserName); err == redis.Nil {
		if err := u.jwtDao.SetRedisJWT(ctx, token, user.UserName); err != nil {
			return systemRes.LoginResponse{}, err
		}

		global.MPA_DB.Model(&user).UpdateColumn("last_login", time.Now().Local())

		return systemRes.LoginResponse{
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, nil

	} else if err != nil {
		return systemRes.LoginResponse{}, errors.WithCode(errcode.ErrTokenRedisStateSetFailed, err.Error())
	} else {
		var blackJWT system.JwtBlacklist
		blackJWT.Jwt = jwtStr

		if err := u.jwtDao.JsonInBlacklist(ctx, blackJWT); err != nil {
			return systemRes.LoginResponse{}, errors.WithCode(errcode.ErrExpired, "%s", err.Error())
		}

		if err := u.jwtDao.SetRedisJWT(ctx, token, user.UserName); err != nil {
			return systemRes.LoginResponse{}, err
		}

		u.userDao.UpdateLoginTime(ctx, user.ID)
		return systemRes.LoginResponse{
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, nil
	}
}

// ChangePassword 修改密码
// 参数：
//
//	userID: 用户 ID
//	req: 修改密码请求信息
//
// 返回值：
//
//	error: 错误信息
func (u *UserService) ChangePassword(ctx context.Context, userID uint, req systemReq.ChangePasswordReq) error {
	user := system.SysUser{MPA_MODEL: global.MPA_MODEL{ID: userID}, Password: req.Password}

	_, err := u.userDao.ChangePassword(ctx, user, req.NewPassword)
	return err
}

// GetUserList 获取用户列表
// 参数：
//
//	authorityID: 权限 ID
//	pageInfo: 分页信息
//
// 返回值：
//
//	response.PageResult: 分页结果
//	error: 错误信息
func (u *UserService) GetUserList(ctx context.Context, authorityID uint, pageInfo systemReq.GetUserList) (response.PageResult, error) {
	list, total, err := u.userDao.GetUserInfoList(ctx, pageInfo, authorityID)
	if err != nil {
		return response.PageResult{}, err
	}

	return response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, nil
}

// SetUserAuthority 设置用户权限
// 参数：
//
//	userID: 用户 ID
//	authorityID: 权限 ID
//	claims: JWT 声明信息
//
// 返回值：
//
//	string: 新的 token
//	int64: 过期时间
//	error: 错误信息
func (u *UserService) SetUserAuthority(ctx context.Context, userID uint, authorityID uint, claims *systemReq.CustomClaims) (string, int64, error) {
	if err := u.userDao.SetUserAuthority(ctx, userID, authorityID); err != nil {
		return "", 0, err
	}

	j := &utils.JWT{SigningKey: []byte(global.MPA_CONFIG.JWT.SigningKey)}
	claims.AuthorityId = authorityID

	token, err := j.CreateToken(*claims)
	if err != nil {
		return "", 0, err
	}

	return token, claims.ExpiresAt.Unix(), nil
}

// SetUserAuthorities 设置用户权限列表
// 参数：
//
//	userID: 用户 ID
//	authorityIDs: 权限 ID 列表
//
// 返回值：
//
//	error: 错误信息
func (u *UserService) SetUserAuthorities(ctx context.Context, userID uint, authorityIDs []uint) error {
	return u.userDao.SetUserAuthorities(ctx, userID, authorityIDs)
}

// DeleteUser 删除用户
// 参数：
//
//	userID: 要删除的用户 ID
//	jwtID: 当前登录用户的 ID
//
// 返回值：
//
//	error: 错误信息
func (u *UserService) DeleteUser(ctx context.Context, userID uint, jwtID uint) error {
	if jwtID == userID {
		return errors.WithCode(errcode.ErrUserDeleteSelf, "禁止删除")
	}

	return u.userDao.DeleteUser(ctx, int(userID))
}

// SetUserInfo 设置用户信息
// 参数：
//
//	user: 用户信息
//
// 返回值：
//
//	error: 错误信息
func (u *UserService) SetUserInfo(ctx context.Context, user systemReq.ChangeUserInfo) error {
	if len(user.AuthorityIds) != 0 {
		if err := u.userDao.SetUserAuthorities(ctx, user.ID, user.AuthorityIds); err != nil {
			return err
		}
	}

	return u.userDao.SetUserInfo(ctx, system.SysUser{
		MPA_MODEL: global.MPA_MODEL{
			ID: user.ID,
		},
		NickName: user.NickName,
		Avatar:   user.Avatar,
		Phone:    user.Phone,
		Email:    user.Email,
		IsActive: user.IsActive,
	})
}

// SetSelfInfo 设置当前用户信息
// 参数：
//
//	userID: 用户 ID
//	user: 用户信息
//
// 返回值：
//
//	error: 错误信息
func (u *UserService) SetSelfInfo(ctx context.Context, userID uint, user systemReq.ChangeUserInfo) error {
	user.ID = userID
	return u.userDao.SetSelfInfo(ctx, system.SysUser{
		MPA_MODEL: global.MPA_MODEL{
			ID: user.ID,
		},
		NickName: user.NickName,
		Avatar:   user.Avatar,
		Phone:    user.Phone,
		Email:    user.Email,
		IsActive: user.IsActive,
	})
}

// GetUserInfo 获取用户信息
// 参数：
//
//	uuid: 用户 UUID
//
// 返回值：
//
//	system.SysUser: 用户信息
//	error: 错误信息
func (u *UserService) GetUserInfo(ctx context.Context, uuid uuid.UUID) (system.SysUser, error) {
	return u.userDao.GetUserInfo(ctx, uuid)
}

// ResetPassword 重置用户密码
// 参数：
//
//	userID: 用户 ID
//
// 返回值：
//
//	error: 错误信息
func (u *UserService) ResetPassword(ctx context.Context, userID uint) error {
	return u.userDao.ResetPassword(ctx, userID)
}
