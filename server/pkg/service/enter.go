package service

import (
	"context"
	"server/pkg/api/request"
	"server/pkg/api/response"
	sysResp "server/pkg/model/common/response"
	"server/pkg/model/system"

	uuid "github.com/satori/go.uuid"
)

// 服务接口定义
type CasbinServiceInterface interface {
	SetCasbin(ctx context.Context, adminAuthorityID, authorityID uint, casbinInfos []request.CasbinInfo) error
	GetPolicyPathByAuthorityID(ctx context.Context, authorityID string) []request.CasbinInfo
}

type JwtServiceInterface interface {
	GenerateToken(ctx context.Context, user system.SysUser) (response.LoginResponse, error)
	AddTokenToBlacklist(ctx context.Context, token string) error
}

type UserServiceInterface interface {
	Register(ctx context.Context, req request.Register) (response.SysUserResponse, error)
	Login(ctx context.Context, l request.Login) (response.LoginResponse, error)
	ChangePassword(ctx context.Context, userID uint, req request.ChangePasswordReq) error
	GetUserList(ctx context.Context, authorityID uint, pageInfo request.GetUserList) (sysResp.PageResult, error)
	SetUserAuthority(ctx context.Context, userID uint, authorityID uint, claims *request.CustomClaims) (string, int64, error)
	SetUserAuthorities(ctx context.Context, userID uint, authorityIDs []uint) error
	DeleteUser(ctx context.Context, userID uint, jwtID uint) error
	SetUserInfo(ctx context.Context, user request.ChangeUserInfo) error
	SetSelfInfo(ctx context.Context, userID uint, user request.ChangeUserInfo) error
	GetUserInfo(ctx context.Context, uuid uuid.UUID) (system.SysUser, error)
	ResetPassword(ctx context.Context, userID uint) error
}

type AuthorityServiceInterface interface {
	GetAuthorityList(ctx context.Context, authorityID uint) ([]system.SysRoleAuthority, error)
	AddAuthority(ctx context.Context, authority system.SysRoleAuthority) (system.SysRoleAuthority, error)
	DeleteAuthority(ctx context.Context, authority system.SysRoleAuthority) error
	UpdateAuthority(ctx context.Context, auth system.SysRoleAuthority) (system.SysRoleAuthority, error)
	GetAuthorityMenu(ctx context.Context, authorityID uint) ([]system.SysAuthorizedMenu, error)
	SetMenuAuthority(ctx context.Context, authorityMenu request.SetMenuAuthority, authorityID, adminAuthorityID uint) error
}

type MenuServiceInterface interface {
	GetAsyncMenu(ctx context.Context, authorityID uint) ([]system.SysAuthorizedMenu, error)
	GetMenuList(ctx context.Context, authorityID uint) ([]system.SysMenu, error)
	UpdateMenu(ctx context.Context, menu system.SysMenu) error
	AddMenu(ctx context.Context, menu system.SysMenu) error
	DeleteMenu(ctx context.Context, menuID uint) error
}

type ApiServiceInterface interface {
	CreateAPI(context.Context, system.SysApi) (system.SysApi, error)
	GetAPIList(context.Context, request.SearchApiParams) (sysResp.PageResult, error)
	GetAllAPIs(context.Context, uint) ([]system.SysApi, error)
	DeleteAPI(context.Context, uint) error
	UpdateAPI(context.Context, system.SysApi) (system.SysApi, error)
	GetAPIGroups(context.Context) ([]string, error)
	SyncAPI(context.Context) ([]system.SysApi, []system.SysApi, []system.SysApi, error)
	IgnoreAPI(context.Context, system.SysIgnoreApi) error
	EnterSyncAPI(context.Context, request.EnterSyncApiParams) error
}
