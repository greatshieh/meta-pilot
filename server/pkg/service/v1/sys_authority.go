package system

import (
	"context"
	"server/pkg/api/request"
	"server/pkg/errcode"
	"server/pkg/model/system"
	"strconv"

	"github.com/marmotedu/errors"
)

// GetAuthorityList 获取权限列表
// 参数：
//
//	authorityID: 权限 ID
//
// 返回值：
//
//	[]system.SysRoleAuthority: 权限列表
//	error: 错误信息
//
// 功能：
//
//	根据用户权限返回相应的权限数据
func (a *AuthorityService) GetAuthorityList(ctx context.Context, authorityID uint) ([]system.SysRoleAuthority, error) {
	return a.authorityDao.GetAuthorityList(ctx, authorityID)
}

// AddAuthority 添加权限
// 参数：
//
//	authority: 要添加的权限信息
//
// 返回值：
//
//	system.SysRoleAuthority: 添加的权限信息
//	error: 错误信息
//
// 功能：
//  1. 接收前端传递的权限信息
//  2. 验证请求参数的合法性
//  3. 处理权限父子关系（在严格权限模式下）
//  4. 调用服务层方法创建权限
//  5. 返回操作结果
func (a *AuthorityService) AddAuthority(ctx context.Context, authority system.SysRoleAuthority) (system.SysRoleAuthority, error) {
	var rules = make([][]string, 0)
	err := a.tx.Transaction(func(ctx context.Context) error {
		// 调用服务层创建权限
		defaultCabinInofs, err := a.authorityDao.AddAuthority(ctx, authority)
		if err != nil {
			return err
		}

		authorityIDStr := strconv.Itoa(int(authority.AuthorityID))

		for _, v := range defaultCabinInofs {
			rules = append(rules, []string{authorityIDStr, v.Path, v.Method})
		}

		return a.casbinDao.AddPolicies(ctx, rules)
	})

	if err != nil {
		return authority, errors.WithCode(errcode.ErrAuthorityCreateFailed, err.Error())
	}

	return authority, nil
}

// DeleteAuthority 删除权限
// 参数：
//
//	authority: 要删除的权限信息
//
// 返回值：
//
//	error: 错误信息
func (a *AuthorityService) DeleteAuthority(ctx context.Context, authority system.SysRoleAuthority) error {
	err := a.tx.Transaction(func(ctx context.Context) error {
		if err := a.authorityDao.DeleteAuthority(ctx, authority); err != nil {
			return err
		}

		// 下面是被注释掉的Casbin权限策略删除代码
		// 如果系统使用Casbin进行权限控制，可以取消注释并使用下面的代码删除相关策略
		authorityIDStr := strconv.Itoa(int(authority.AuthorityID))
		if err := a.casbinDao.RemoveFilteredPolicy(ctx, authorityIDStr); err != nil {
			return err
		}

		// 事务执行成功，返回nil
		return nil
	})

	if err != nil {
		return errors.WithCode(errcode.ErrAuthorityDeleteFailed, err.Error())
	}

	return nil

}

// UpdateAuthority 更新权限
// 参数：
//
//	auth: 要更新的权限信息
//
// 返回值：
//
//	system.SysRoleAuthority: 更新后的权限信息
//	error: 错误信息
func (a *AuthorityService) UpdateAuthority(ctx context.Context, auth system.SysRoleAuthority) (system.SysRoleAuthority, error) {
	return a.authorityDao.UpdateAuthority(ctx, auth)
}

// GetAuthorityMenu 获取权限菜单
// 参数：
//
//	authorityID: 权限 ID
//
// 返回值：
//
//	[]system.SysAuthorizedMenu: 权限菜单列表
//	error: 错误信息
func (a *AuthorityService) GetAuthorityMenu(ctx context.Context, authorityID uint) ([]system.SysAuthorizedMenu, error) {
	return a.menuDao.GetMenuAuthority(ctx, authorityID)
}

// SetMenuAuthority 设置菜单权限
// 参数：
//
//	authorityMenu: 菜单权限设置
//	authorityID: 权限 ID
//	adminAuthorityID: 管理员权限 ID
//
// 返回值：
//
//	error: 错误信息
func (a *AuthorityService) SetMenuAuthority(ctx context.Context, authorityMenu request.SetMenuAuthority, authorityID uint, adminAuthorityID uint) error {
	return a.authorityDao.SetMenuAuthority(ctx, authorityMenu.Menus, adminAuthorityID, authorityID)
}
