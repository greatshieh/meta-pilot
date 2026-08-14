// 统一DAO接口
package dao

import (
	"context"
	"server/pkg/api/request"
	sysReq "server/pkg/model/common/request"
	"server/pkg/model/system"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type BaseDaoInterface interface {
	GetDB(ctx context.Context) *gorm.DB
}

// BaseDao 基础DAO结构体
type BaseDao struct {
	db *gorm.DB
}

// NewBaseDao 创建基础DAO实例
func NewBaseDao(db *gorm.DB) *BaseDao {
	return &BaseDao{
		db,
	}
}

// GetDB 获取数据库连接
func (d *BaseDao) GetDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value("tx").(*gorm.DB); ok {
		return tx
	}
	return d.db.WithContext(ctx)
}

// UserDaoInterface 用户DAO接口
// 提供用户相关的数据访问方法
type UserDaoInterface interface {
	BaseDaoInterface // 强制实现GetDB方法
	Register(context.Context, system.SysUser) (system.SysUser, error)
	Login(context.Context, system.SysUser) (*system.SysUser, error)
	ChangePassword(context.Context, system.SysUser, string) (*system.SysUser, error)
	GetUserInfoList(context.Context, request.GetUserList, uint) ([]system.SysUser, int64, error)
	SetUserAuthority(ctx context.Context, userID, authorityID uint) error
	SetUserAuthorities(context.Context, uint, []uint) error
	DeleteUser(context.Context, int) error
	SetUserInfo(context.Context, system.SysUser) error
	SetSelfInfo(context.Context, system.SysUser) error
	GetUserInfo(context.Context, uuid.UUID) (system.SysUser, error)
	FindUserByID(context.Context, int) (*system.SysUser, error)
	FindUserByUUID(context.Context, string) (*system.SysUser, error)
	ResetPassword(context.Context, uint) error
	UpdateLoginTime(context.Context, uint) error
}

// JwtDaoInterface JWT DAO接口
// 提供JWT相关的数据访问方法
type JwtDaoInterface interface {
	BaseDaoInterface // 强制实现GetDB方法
	JsonInBlacklist(context.Context, system.JwtBlacklist) error
	GetRedisJWT(context.Context, string) (string, error)
	SetRedisJWT(ctx context.Context, string, userName string) error
	IsBlacklist(context.Context, string) bool
}

// CasbinDaoInterface Casbin DAO接口
// 提供Casbin相关的数据访问方法
type CasbinDaoInterface interface {
	BaseDaoInterface // 强制实现GetDB方法
	SetCasbin(context.Context, [][]string) error
	UpdateCasbinApi(ctx context.Context, oldPath, newPath, oldMethod, newMethod string) error
	GetPolicyPathByAuthorityID(context.Context, string) []request.CasbinInfo
	ClearCasbin(context.Context, int, ...string) (bool, error)
	RemoveFilteredPolicy(context.Context, string) error
	SyncPolicy(context.Context, string, [][]string) error
	AddPolicies(context.Context, [][]string) error
	FreshCasbin(context.Context) (err error)
}

// AuthorityDaoInterface 权限DAO接口
// 提供权限相关的数据访问方法
type AuthorityDaoInterface interface {
	BaseDaoInterface // 强制实现GetDB方法
	GetAuthorityList(context.Context, uint) ([]system.SysRoleAuthority, error)
	DeleteAuthority(context.Context, system.SysRoleAuthority) error
	UpdateAuthority(context.Context, system.SysRoleAuthority) (system.SysRoleAuthority, error)
	SetMenuAuthority(ctx context.Context, menus []system.SysMenu, adminAuthorityID, authorityID uint) (err error)
	GetParentAuthorityID(context.Context, uint) (uint, error)
	CheckAuthorityIDAuth(ctx context.Context, authorityID, targetID uint) error
	AddAuthority(context.Context, system.SysRoleAuthority) ([]system.SysApi, error)
}

// MenuDaoInterface 菜单DAO接口
// 提供菜单相关的数据访问方法
type MenuDaoInterface interface {
	BaseDaoInterface // 强制实现GetDB方法
	GetAsyncMenu(context.Context, uint) ([]system.SysAuthorizedMenu, error)
	UpdateMenu(context.Context, system.SysMenu) error
	AddMenu(context.Context, system.SysMenu) error
	DeleteMenu(context.Context, uint) error
	GetMenuList(ctx context.Context, parentAuthorityID, authorityID uint) ([]system.SysMenu, error)
	GetMenuAuthority(context.Context, uint) ([]system.SysAuthorizedMenu, error)
}

// ApiDaoInterface API DAO接口
// 提供API相关的数据访问方法
type ApiDaoInterface interface {
	BaseDaoInterface // 强制实现GetDB方法
	CreateAPI(context.Context, system.SysApi) error
	DeleteAPI(context.Context, uint) (system.SysApi, error)
	GetAPIInfoList(context.Context, system.SysApi, sysReq.PageInfo) ([]system.SysApi, int64, error)
	UpdateAPI(context.Context, system.SysApi) error
	GetAllAPIs(context.Context, uint) ([]system.SysApi, error)
	GetAPIGroups(context.Context) ([]string, error)
	SyncAPI(context.Context) (newApis, deleteApis, ignoreApis []system.SysApi, err error)
	IgnoreAPI(context.Context, system.SysIgnoreApi) error
	EnterSyncAPI(context.Context, request.EnterSyncApiParams) error
}
