package system

import (
	"context"
	"server/pkg/model/system"
)

// GetAsyncMenu 获取异步菜单
// 参数：
//
//	authorityID: 权限 ID
//
// 返回值：
//
//	[]system.SysAuthorizedMenu: 菜单列表
//	error: 错误信息
//
// 功能：
//
//	根据用户权限返回相应的菜单数据
func (m *MenuService) GetAsyncMenu(ctx context.Context, authorityID uint) ([]system.SysAuthorizedMenu, error) {
	// 调用菜单DAO获取异步菜单数据，传入用户的权限ID作为参数
	// 这样可以确保用户只能看到自己有权限访问的菜单项
	menus, err := m.menuDao.GetAsyncMenu(ctx, authorityID)

	// 检查返回的菜单数据是否为nil（空）
	// 如果为空，则初始化为一个空的菜单切片
	if menus == nil {
		menus = []system.SysAuthorizedMenu{}
	}

	return menus, err
}

// GetMenuList 获取菜单列表
// 参数：
//
//	authorityID: 权限 ID
//
// 返回值：
//
//	[]system.SysMenu: 菜单列表
//	error: 错误信息
//
// 功能：
// 主要流程：
//  1. 调用菜单服务获取该权限ID对应的菜单列表
//  2. 处理可能的错误情况
//  3. 确保返回空数组而不是null
//  4. 将结果以JSON格式返回给前端
//     根据用户的权限ID获取对应的菜单列表
func (m *MenuService) GetMenuList(ctx context.Context, authorityID uint) ([]system.SysMenu, error) {
	parentAuthorityID, err := m.authorityDao.GetParentAuthorityID(ctx, authorityID)
	if err != nil {
		return nil, err
	}

	// 调用菜单DAO的GetMenuList方法，根据权限ID获取菜单列表
	// GetMenuList会构建完整的菜单树结构，包含所有层级的子菜单
	menus, err := m.menuDao.GetMenuList(ctx, parentAuthorityID, authorityID)
	if err != nil {
		return nil, err
	}

	// 检查返回的菜单列表是否为nil（防止返回null值）
	// 如果为nil，则初始化为空的菜单数组
	if menus == nil {
		menus = []system.SysMenu{}
	}

	return menus, err
}

// UpdateMenu 更新菜单信息
// 参数：
//
//	menu: 菜单信息
//
// 返回值：
//
//	error: 错误信息
//
// 功能：
//
//	更新菜单信息
func (m *MenuService) UpdateMenu(ctx context.Context, menu system.SysMenu) error {
	// 调用菜单DAO的UpdateMenu方法执行菜单更新操作
	return m.menuDao.UpdateMenu(ctx, menu)
}

// AddMenu 添加新菜单
// 参数：
//
//	menu: 菜单信息
//
// 返回值：
//
//	error: 错误信息
//
// 功能：
//
//	添加新菜单
func (m *MenuService) AddMenu(ctx context.Context, menu system.SysMenu) error {
	// 调用菜单DAO的AddMenu方法执行菜单创建操作
	return m.menuDao.AddMenu(ctx, menu)
}

// DeleteMenu 删除菜单
// 参数：
//
//	menuID: 菜单 ID
//
// 返回值：
//
//	error: 错误信息
//
// 功能：
//
//	删除指定ID的菜单
func (m *MenuService) DeleteMenu(ctx context.Context, menuID uint) error {
	// 调用菜单DAO的DeleteMenu方法执行菜单删除操作
	return m.menuDao.DeleteMenu(ctx, menuID)
}
