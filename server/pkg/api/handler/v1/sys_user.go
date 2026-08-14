package v1

import (
	"server/pkg/errcode"
	"server/pkg/model/common/response"
	"server/pkg/service"
	"server/pkg/utils"
	"strconv"

	"server/pkg/model/system"

	"server/pkg/api/handler"
	systemReq "server/pkg/api/request"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
)

// UserHandler 用户处理器结构体
// 处理用户相关的HTTP请求
type UserHandler struct {
	userService service.UserServiceInterface
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(userService service.UserServiceInterface) handler.UserHandlerInterface {
	return &UserHandler{
		userService: userService,
	}
}

// Register 用户注册
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理用户注册请求，验证请求参数，调用服务层进行注册，返回注册结果
func (u *UserHandler) Register(c *gin.Context) {
	var r systemReq.Register

	if err := utils.VerifyBindJson(c, &r, utils.RegisterVerify); err != nil {
		response.FailResponse(c, err)
		return
	}

	userReturn, err := u.userService.Register(c.Request.Context(), r)

	if err != nil {
		response.FailResponse(c, err)
		return
	}

	response.SuccessResponse(c, userReturn, "注册成功")
}

// Login 用户登录
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理用户登录请求，验证请求参数，调用服务层进行登录，返回登录结果和token
func (u *UserHandler) Login(c *gin.Context) {
	var l systemReq.Login

	if err := utils.VerifyBindJson(c, &l, utils.LoginVerify); err != nil {
		response.FailResponse(c, err)
		return
	}

	userInter, err := u.userService.Login(c.Request.Context(), l)

	if err != nil {
		response.FailResponse(c, err)
		return
	}

	response.SuccessResponse(c, userInter, "登录成功")
}

// ChangePassword 修改密码
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理用户修改密码请求，验证请求参数，调用服务层修改密码，返回修改结果
func (u *UserHandler) ChangePassword(c *gin.Context) {
	var req systemReq.ChangePasswordReq

	if err := utils.VerifyBindJson(c, &req, utils.ChangePasswordVerify); err != nil {
		response.FailResponse(c, err)
		return
	}

	id := utils.GetUserID(c)

	if err := u.userService.ChangePassword(c.Request.Context(), id, req); err != nil {
		response.FailResponse(c, err)
		return
	}

	response.SuccessResponse(c, nil, "密码修改成功")
}

// GetUserList 获取用户列表
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理获取用户列表请求，验证请求参数，调用服务层获取用户列表，返回分页结果
func (u *UserHandler) GetUserList(c *gin.Context) {
	var pageInfo systemReq.GetUserList

	if err := utils.VerifyBindQuery(c, &pageInfo, utils.PageInfoVerify); err != nil {
		response.FailResponse(c, err)
		return
	}

	authorityID := utils.GetUserAuthorityId(c)

	pageResult, err := u.userService.GetUserList(c.Request.Context(), authorityID, pageInfo)
	if err != nil {
		response.FailResponse(c, err)
		return
	}

	response.SuccessResponse(c, pageResult)
}

// SetUserAuthority 设置用户权限
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理设置用户权限请求，验证请求参数，调用服务层设置用户权限，返回设置结果和新token
func (u *UserHandler) SetUserAuthority(c *gin.Context) {
	var sua systemReq.SetUserAuth

	if err := utils.VerifyBindJson(c, &sua, utils.SetUserAuthorityVerify); err != nil {
		response.FailResponse(c, err)
		return
	}

	userID := utils.GetUserID(c)
	claims := utils.GetUserInfo(c)

	token, expiresAtTime, err := u.userService.SetUserAuthority(c.Request.Context(), userID, sua.AuthorityId, claims)

	if err != nil {
		response.FailResponse(c, err)
		return
	}

	c.Header("new-token", token)
	c.Header("new-expires-at", strconv.FormatInt(expiresAtTime, 10))
	response.SuccessResponse(c, nil, "切换权限成功")
}

// SetUserAuthorities 设置用户权限列表
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理设置用户权限列表请求，验证请求参数，调用服务层设置用户权限列表，返回设置结果
func (u *UserHandler) SetUserAuthorities(c *gin.Context) {
	var sua systemReq.SetUserAuthorities

	if err := utils.VerifyBindJson(c, &sua, utils.SetUserAuthorityVerifies); err != nil {
		response.FailResponse(c, err)
		return
	}

	if err := u.userService.SetUserAuthorities(c.Request.Context(), sua.ID, sua.AuthorityIds); err != nil {
		response.FailResponse(c, err)
		return
	}

	response.SuccessResponse(c, nil, "权限设置成功")
}

// DeleteUser 删除用户
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理删除用户请求，验证请求参数，调用服务层删除用户，返回删除结果
func (u *UserHandler) DeleteUser(c *gin.Context) {
	reqID := c.Param("id")

	id, _ := strconv.Atoi(reqID)

	jwtID := utils.GetUserID(c)

	if jwtID == uint(id) {
		response.FailResponse(c, errors.WithCode(errcode.ErrUserDeleteSelf, "禁止删除"))
		return
	}

	if err := u.userService.DeleteUser(c.Request.Context(), uint(id), jwtID); err != nil {
		response.FailResponse(c, err)
		return
	}

	response.SuccessResponse(c, nil, "删除成功")
}

// SetUserInfo 设置用户信息
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理设置用户信息请求，验证请求参数，调用服务层设置用户信息，返回设置结果
func (u *UserHandler) SetUserInfo(c *gin.Context) {
	var user systemReq.ChangeUserInfo

	if err := utils.VerifyBindJson(c, &user, utils.IdVerify); err != nil {
		response.FailResponse(c, err)
		return
	}

	if len(user.AuthorityIds) != 0 {
		if err := u.userService.SetUserAuthorities(c.Request.Context(), user.ID, user.AuthorityIds); err != nil {
			response.FailResponse(c, err)
			return
		}
	}

	err := u.userService.SetUserInfo(c.Request.Context(), user)
	if err != nil {
		response.FailResponse(c, err)
		return
	}

	response.SuccessResponse(c, nil, "设置员工信息成功")
}

// SetSelfInfo 设置当前用户信息
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理设置当前用户信息请求，验证请求参数，调用服务层设置当前用户信息，返回设置结果
func (u *UserHandler) SetSelfInfo(c *gin.Context) {
	var user systemReq.ChangeUserInfo

	if err := utils.VerifyBindJson(c, &user); err != nil {
		response.FailResponse(c, err)
		return
	}

	user.ID = utils.GetUserID(c)

	err := u.userService.SetSelfInfo(c.Request.Context(), user.ID, user)

	if err != nil {
		response.FailResponse(c, err)
		return
	}

	response.SuccessResponse(c, nil, "更新信息成功")
}

// GetUserInfo 获取用户信息
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理获取用户信息请求，调用服务层获取用户信息，返回用户信息
func (u *UserHandler) GetUserInfo(c *gin.Context) {
	uuid := utils.GetUserUuid(c)

	ReqUser, err := u.userService.GetUserInfo(c.Request.Context(), uuid)
	if err != nil {
		response.FailResponse(c, err)
		return
	}

	response.SuccessResponse(c, ReqUser)
}

// ResetPassword 重置用户密码
// 参数：
//
//	c: gin上下文，用于处理HTTP请求和响应
//
// 功能：
//
//	处理重置用户密码请求，验证请求参数，调用服务层重置用户密码，返回重置结果
func (u *UserHandler) ResetPassword(c *gin.Context) {
	var user system.SysUser

	if err := utils.VerifyBindJson(c, &user, utils.IdVerify); err != nil {
		response.FailResponse(c, err)
		return
	}

	if err := u.userService.ResetPassword(c.Request.Context(), user.ID); err != nil {
		response.FailResponse(c, err)
		return
	}

	response.SuccessResponse(c, nil)
}
