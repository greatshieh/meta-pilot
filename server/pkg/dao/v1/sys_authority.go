package v1

import (
	"context"
	"server/pkg/api/request"
	"server/pkg/errcode"
	"server/pkg/global"
	"server/pkg/model/system"
	"strconv"

	"github.com/marmotedu/errors"
	"gorm.io/gorm"
)

// GetAuthorityList 获取权限列表
// 参数：
//
//	authorityID: 权限ID
//
// 返回值：
//
//	list: 权限列表
//	err: 错误信息
//
// 功能：
//
//	根据传入的authorityID获取权限列表，支持严格的树形结构权限控制
//	如果开启了严格权限控制：
//	  - 顶级角色可以查看自己及下属角色
//	  - 非顶级角色只能查看自己的直接下级角色
//
//	如果未开启严格权限控制：
//	  - 可以查看所有顶级角色(父ID为0的角色)
func (authorityDao *AuthorityDao) GetAuthorityList(ctx context.Context, authorityID uint) (list []system.SysRoleAuthority, err error) {
	var authority system.SysRoleAuthority
	err = authorityDao.GetDB(ctx).Where("authority_id = ?", authorityID).First(&authority).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(errcode.ErrAuthorityNotFound, "%s", err.Error())
		}
		return nil, errors.WithCode(errcode.ErrDatabase, "%s", err.Error())
	}

	var authorities []system.SysRoleAuthority
	db := authorityDao.GetDB(ctx).Model(&system.SysRoleAuthority{}).Preload("SysMenu").Preload("Users")
	if global.MPA_CONFIG.System.UseStrictAuth {
		// 当开启了严格树形结构后
		if *authority.ParentID == 0 {
			// 只有顶级角色可以修改自己的权限和以下权限
			err = db.Preload("AuthorizedSubRoles").Where("authority_id = ?", authorityID).Find(&authorities).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = errors.WithCode(errcode.ErrAuthorityNotFound, "%s", err.Error())
			}
		} else {
			// 非顶级角色只能修改以下权限
			err = db.Preload("AuthorizedSubRoles").Where("parent_id = ?", authorityID).Find(&authorities).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = errors.WithCode(errcode.ErrAuthorityNotFound, "%s", err.Error())
			}
		}
	} else {
		err = db.Preload("AuthorizedSubRoles").Where("parent_id = ?", "0").Find(&authorities).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.WithCode(errcode.ErrAuthorityNotFound, "%s", err.Error())
		}
	}

	// 递归查找子权限和用户
	for k := range authorities {
		err = authorityDao.findChildrenAuthority(ctx, &authorities[k])
	}

	return authorities, err
}

// findChildrenAuthority 递归查找子权限
// 参数：
//
//	authority: 权限对象指针
//
// 返回值：
//
//	err: 错误信息
//
// 功能：
//
//	通过递归方式查找指定权限的所有子权限，并加载其数据权限
func (authorityDao *AuthorityDao) findChildrenAuthority(ctx context.Context, authority *system.SysRoleAuthority) (err error) {
	err = authorityDao.GetDB(ctx).Preload("AuthorizedSubRoles").Preload("Users").Where("parent_id = ?", authority.AuthorityID).Find(&authority.Children).Error
	if err != nil {
		return errors.WithCode(errcode.ErrDatabase, "%s", err.Error())
	}

	if len(authority.Children) > 0 {
		for k := range authority.Children {
			if err = authorityDao.findChildrenAuthority(ctx, &authority.Children[k]); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetParentAuthorityID 获取父级权限ID
// 参数：
//
//	authorityID: 权限ID
//
// 返回值：
//
//	parentID: 父级权限ID
//	err: 错误信息
//
// 功能：
//
//	根据传入的权限ID查询其父级权限ID
func (authorityDao *AuthorityDao) GetParentAuthorityID(ctx context.Context, authorityID uint) (parentID uint, err error) {
	var authority system.SysRoleAuthority
	err = authorityDao.GetDB(ctx).Where("authority_id = ?", authorityID).First(&authority).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return parentID, errors.WithCode(errcode.ErrAuthorityNotFound, "%s", err.Error())
		}
		return parentID, errors.WithCode(errcode.ErrDatabase, "%s", err.Error())
	}

	return *authority.ParentID, nil
}

// getStructAuthorityList 获取权限结构列表
// 参数：
//
//	authorityID: 权限ID
//
// 返回值：
//
//	list: 权限ID列表
//	err: 错误信息
//
// 功能：
//
//	递归获取指定权限ID下的所有子权限ID列表
//	如果是顶级权限，则将自身也加入到列表中
func (authorityDao *AuthorityDao) getStructAuthorityList(ctx context.Context, authorityID uint) (list []uint, err error) {
	var auth system.SysRoleAuthority
	_ = authorityDao.GetDB(ctx).First(&auth, "authority_id = ?", authorityID).Error
	var authorities []system.SysRoleAuthority
	err = authorityDao.GetDB(ctx).Preload("AuthorizedSubRoles").Where("parent_id = ?", authorityID).Find(&authorities).Error
	if len(authorities) > 0 {
		for k := range authorities {
			list = append(list, authorities[k].AuthorityID)
			childrenList, err := authorityDao.getStructAuthorityList(ctx, authorities[k].AuthorityID)
			if err == nil {
				list = append(list, childrenList...)
			}
		}
	}
	if *auth.ParentID == 0 {
		list = append(list, authorityID)
	}
	return list, err
}

// CheckAuthorityIDAuth 检查权限ID合法性
// 参数：
//
//	authorityID: 当前用户权限ID
//	targetID: 目标权限ID
//
// 返回值：
//
//	err: 错误信息
//
// 功能：
//
//	在开启严格权限验证时，检查目标权限ID是否在当前用户权限范围内
func (authorityDao *AuthorityDao) CheckAuthorityIDAuth(ctx context.Context, authorityID, targetID uint) (err error) {
	if !global.MPA_CONFIG.System.UseStrictAuth {
		return nil
	}
	authIDS, err := authorityDao.getStructAuthorityList(ctx, authorityID)
	if err != nil {
		return err
	}
	hasAuth := false
	for _, v := range authIDS {
		if v == targetID {
			hasAuth = true
			break
		}
	}
	if !hasAuth {
		return errors.New("您提交的角色ID不合法")
	}
	return nil
}

// AddAuthority 添加权限
// 参数：
//
//	authority: 权限信息
//
// 返回值：
//
//	authority: 创建成功的权限信息
//	err: 错误信息
//
// 功能：
//
//	创建一个新的权限角色，确保角色ID唯一性，并为其分配默认菜单权限
func (authorityDao *AuthorityDao) AddAuthority(ctx context.Context, authority system.SysRoleAuthority) ([]system.SysApi, error) {
	// 设置默认api权限
	var defaultCabinInofs []system.SysApi

	// 检查是否已存在相同的权限ID，防止重复创建
	// 使用gorm的First方法查询是否存在相同authority_id的记录
	// 如果查询结果不是"记录未找到"错误，说明已存在相同ID的权限
	if err := authorityDao.GetDB(ctx).Where("authority_id = ?", authority.AuthorityID).First(&system.SysRoleAuthority{}).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return defaultCabinInofs, errors.WithCode(errcode.ErrAuthorityAlreadyExist, "存在相同角色id")
	}

	if err := authorityDao.GetDB(ctx).Create(&authority).Error; err != nil {
		// 创建失败，返回错误
		return defaultCabinInofs, err
	}

	// 为新创建的权限分配默认菜单权限
	// 获取系统默认菜单列表
	authority.SysMenu = request.DefaultMenu()
	// 使用GORM的Association功能替换权限与菜单的关联关系
	// 这会清空原有的关联并建立新的关联
	if err := authorityDao.GetDB(ctx).Model(&authority).Association("SysMenu").Replace(&authority.SysMenu); err != nil {
		// 关联失败，返回错误
		return defaultCabinInofs, err
	}

	if err := authorityDao.GetDB(ctx).Where("required = true").Find(&defaultCabinInofs).Error; err != nil {
		return defaultCabinInofs, err
	}

	return defaultCabinInofs, nil
}

// DeleteAuthority 删除权限
// 参数：
//
//	auth: 权限对象指针
//
// 返回值：
//
//	err: 错误信息
//
// 功能：
//
//	删除指定的权限角色，在删除前进行多项安全检查以确保数据一致性
func (authorityDao *AuthorityDao) DeleteAuthority(ctx context.Context, auth system.SysRoleAuthority) error {
	if errors.Is(authorityDao.GetDB(ctx).Preload("Users").First(&auth).Error, gorm.ErrRecordNotFound) {
		return errors.WithCode(errcode.ErrAuthorityNotFound, "该角色不存在")
	}

	// 检查是否有用户正在使用该权限
	// 如果Users关联数据不为空，说明有用户使用此权限，禁止删除
	if len(auth.Users) != 0 {
		return errors.WithCode(errcode.ErrAuthorityAlreadyAssignedToUser, "此角色有用户正在使用禁止删除")
	}

	if !errors.Is(authorityDao.GetDB(ctx).Where("authority_id = ?", auth.AuthorityID).First(&system.SysUser{}).Error, gorm.ErrRecordNotFound) {
		return errors.WithCode(errcode.ErrAuthorityAlreadyAssignedToUser, "此角色有用户正在使用禁止删除")
	}

	if !errors.Is(authorityDao.GetDB(ctx).Where("parent_id = ?", auth.AuthorityID).First(&system.SysRoleAuthority{}).Error, gorm.ErrRecordNotFound) {
		return errors.WithCode(errcode.ErrAuthorityHasChildren, "此角色存在子角色不允许删除")
	}

	if err := authorityDao.GetDB(ctx).Preload("SysMenu").Preload("AuthorizedSubRoles").Where("authority_id = ?", auth.AuthorityID).First(&auth).Error; err != nil {
		return err
	}

	if err := authorityDao.GetDB(ctx).Unscoped().Delete(&auth).Error; err != nil {
		return err
	}

	if len(auth.SysMenu) > 0 {
		if err := authorityDao.GetDB(ctx).Model(&auth).Association("SysMenu").Delete(auth.SysMenu); err != nil {
			return err
		}
	}

	if len(auth.AuthorizedSubRoles) > 0 {
		if err := authorityDao.GetDB(ctx).Model(&auth).Association("AuthorizedSubRoles").Delete(auth.AuthorizedSubRoles); err != nil {
			return err
		}
	}

	if err := authorityDao.GetDB(ctx).Delete(&system.SysUserAuthorityRelation{}, "sys_authority_authority_id = ?", auth.AuthorityID).Error; err != nil {
		return err
	}

	return nil
}

// UpdateAuthority 更新权限
// 参数：
//
//	auth: 权限信息
//
// 返回值：
//
//	authority: 更新后的权限信息
//	err: 错误信息
//
// 功能：
//
//	更新指定的权限信息
func (authorityDao *AuthorityDao) UpdateAuthority(ctx context.Context, auth system.SysRoleAuthority) (authority system.SysRoleAuthority, err error) {
	var oldAuthority system.SysRoleAuthority
	err = authorityDao.GetDB(ctx).Where("authority_id = ?", auth.AuthorityID).First(&oldAuthority).Error
	if err != nil {
		return system.SysRoleAuthority{}, errors.WithCode(errcode.ErrAuthoritySearchFailed, "查询角色数据失败")
	}
	err = authorityDao.GetDB(ctx).Model(&oldAuthority).Updates(&auth).Error
	if err != nil {
		return auth, errors.WithCode(errcode.ErrAuthorityUpdateFailed, "%s", err.Error())
	}

	return auth, nil
}

// SetMenuAuthority 设置菜单权限
// 参数：
//
//	menus: 要分配给权限的菜单列表
//	adminAuthorityID: 管理员权限ID（执行操作的用户权限）
//	authorityID: 目标权限ID（要被分配菜单权限的角色）
//
// 返回值：
//
//	err: 错误信息
//
// 功能：
//
//	为指定权限分配菜单权限，确保操作符合权限层级限制
func (authorityDao *AuthorityDao) SetMenuAuthority(ctx context.Context, menus []system.SysMenu, adminAuthorityID, authorityID uint) (err error) {
	// 构造权限对象，准备更新菜单权限
	var auth system.SysRoleAuthority
	auth.AuthorityID = authorityID
	auth.SysMenu = menus

	// 检查管理员是否有权为目标权限ID设置菜单权限
	// 在严格权限模式下会验证权限层级关系
	err = authorityDao.CheckAuthorityIDAuth(ctx, adminAuthorityID, authorityID)
	if err != nil {
		return err
	}

	// 查询管理员权限信息，用于后续权限校验
	var authority system.SysRoleAuthority
	_ = global.MPA_DB.First(&authority, "authority_id = ?", adminAuthorityID).Error
	var menuIDs []string

	// 当开启了严格的树角色并且父角色不为0时需要进行菜单筛选
	// 这是为了防止低级权限用户越级分配菜单权限
	if global.MPA_CONFIG.System.UseStrictAuth && *authority.ParentID != 0 {
		// 查询管理员可操作的菜单ID列表
		var authorityMenus []system.SysMenuAuthorityRelation
		err = global.MPA_DB.Where("sys_authority_authority_id = ?", adminAuthorityID).Find(&authorityMenus).Error
		if err != nil {
			return err
		}
		// 将管理员可操作的菜单ID转换为字符串数组
		for i := range authorityMenus {
			menuIDs = append(menuIDs, authorityMenus[i].MenuID)
		}

		// 验证要分配的菜单是否都在管理员可操作范围内
		for i := range menus {
			hasMenu := false
			for j := range menuIDs {
				idStr := strconv.Itoa(int(menus[i].ID))
				if idStr == menuIDs[j] {
					hasMenu = true
				}
			}
			// 如果发现要分配的菜单不在管理员权限范围内，则返回错误
			if !hasMenu {
				return errors.New("添加失败,请勿跨级操作")
			}
		}
	}

	// 执行菜单权限分配操作
	var s system.SysRoleAuthority
	authorityDao.GetDB(ctx).Preload("SysMenu").First(&s, "authority_id = ?", auth.AuthorityID)

	return authorityDao.GetDB(ctx).Model(&s).Association("SysMenu").Replace(&auth.SysMenu)
}
