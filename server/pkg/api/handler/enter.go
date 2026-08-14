package handler

import "github.com/gin-gonic/gin"

// ApiHandlerInterface API处理器接口
type ApiHandlerInterface interface {
	CreateApi(c *gin.Context)
	GetApiList(c *gin.Context)
	GetAllApis(c *gin.Context)
	DeleteApi(c *gin.Context)
	UpdateApi(c *gin.Context)
	GetApiGroups(c *gin.Context)
	SyncApi(c *gin.Context)
	IgnoreApi(c *gin.Context)
	EnterSyncApi(c *gin.Context)
}

// UserHandlerInterface 用户处理器接口
type UserHandlerInterface interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	ChangePassword(c *gin.Context)
	GetUserList(c *gin.Context)
	SetUserAuthority(c *gin.Context)
	SetUserAuthorities(c *gin.Context)
	DeleteUser(c *gin.Context)
	SetUserInfo(c *gin.Context)
	SetSelfInfo(c *gin.Context)
	GetUserInfo(c *gin.Context)
	ResetPassword(c *gin.Context)
}

// CasbinHandlerInterface Casbin处理器接口
type CasbinHandlerInterface interface {
	SetCasbin(c *gin.Context)
	GetPolicyPathByAuthorityId(c *gin.Context)
}

// MenuHandlerInterface 菜单处理器接口
type MenuHandlerInterface interface {
	GetAsyncMenu(c *gin.Context)
	GetMenuList(c *gin.Context)
	UpdateMenu(c *gin.Context)
	AddMenu(c *gin.Context)
	DeleteMenu(c *gin.Context)
}

// AuthorityHandlerInterface 权限处理器接口
type AuthorityHandlerInterface interface {
	GetAuthorityList(c *gin.Context)
	AddAuthority(c *gin.Context)
	DeleteAuthority(c *gin.Context)
	UpdateAuthority(c *gin.Context)
	GetAuthorityMenu(c *gin.Context)
	SetMenuAuthority(c *gin.Context)
}

// JwtHandlerInterface JWT处理器接口
type JwtHandlerInterface interface {
	JsonInBlacklist(c *gin.Context)
}
