package v1

import (
	"context"
	"server/pkg/errcode"
	"server/pkg/global"
	"server/pkg/model/system"

	"github.com/marmotedu/errors"
	"gorm.io/gorm"
)

// GetAsyncMenu 根据权限ID获取异步菜单列表
// 参数：
//
//	authorityID: 用户权限ID
//
// 返回值：
//
//	menus: 菜单列表
//	err: 错误信息
//
// 功能：
//
//	通过构建菜单树结构来返回用户有权访问的菜单项
func (menuDao *MenuDao) GetAsyncMenu(ctx context.Context, authorityID uint) (menus []system.SysAuthorizedMenu, err error) {
	// 调用 getMenuTree 方法获取以权限ID为基础的完整菜单树结构
	menuTree, err := menuDao.getMenuTree(ctx, authorityID)

	// 如果获取菜单树失败，直接返回错误
	if err != nil {
		return nil, err
	}

	// 从菜单树的第一层获取顶级菜单项列表
	// menuTree[0] 包含所有顶级菜单项
	menus = menuTree[0]

	// 遍历所有顶级菜单项，为每个菜单项填充其子菜单列表
	for i := range menus {
		// 调用 getChildrenList 方法递归获取并设置当前菜单项的子菜单列表
		err = menuDao.getChildrenList(&menus[i], menuTree)

		// 如果获取子菜单列表过程中出现错误，则中断循环并返回错误
		if err != nil {
			return menus, err
		}
	}

	// 返回完整的菜单列表和nil错误表示操作成功
	return menus, err
}

// UpdateMenu 更新菜单信息
// 参数：
//
//	menu: 包含更新信息的菜单对象
//
// 返回值：
//
//	err: 错误信息
//
// 功能：
//
//	通过GORM的Updates方法更新菜单记录，只会更新非零值字段
func (menuDao *MenuDao) UpdateMenu(ctx context.Context, menu system.SysMenu) (err error) {
	err = menuDao.GetDB(ctx).Model(&menu).Updates(&menu).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(errcode.ErrMenuSetFailed, "%s", err.Error())
		}
	}
	return nil
}

// AddMenu 添加新菜单
// 参数：
//
//	menu: 包含新菜单信息的菜单对象
//
// 返回值：
//
//	err: 错误信息
//
// 功能：
//
//	通过GORM的Create方法向数据库中插入一条新的菜单记录
func (menuDao *MenuDao) AddMenu(ctx context.Context, menu system.SysMenu) error {
	return menuDao.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("name = ?", menu.Name).First(&system.SysMenu{}).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(errcode.ErrMenuAlreadyExist, "%s", err.Error())
		}

		if menu.ParentID != 0 {
			// 检查父菜单是否存在
			var parentMenu system.SysMenu
			if err := tx.First(&parentMenu, menu.ParentID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.WithCode(errcode.ErrMenuNotFound, "%s", err.Error())
				}
				return errors.WithCode(errcode.ErrDatabase, "%s", err.Error())
			}

			// 检查父菜单下现有子菜单数量
			var existingChildrenCount int64
			err := tx.Model(&system.SysMenu{}).Where("parent_id = ?", menu.ParentID).Count(&existingChildrenCount).Error
			if err != nil {
				return errors.WithCode(errcode.ErrDatabase, "%s", err.Error())
			}
		}

		if err := tx.Create(&menu).Error; err != nil {
			return errors.WithCode(errcode.ErrDatabase, "%s", err.Error())
		}

		return nil
	})
}

// DeleteMenu 删除菜单
// 参数：
//
//	menuID: 菜单ID
//
// 返回值：
//
//	err: 错误信息
//
// 功能：
//
//	删除指定ID的菜单及其关联数据
func (menuDao *MenuDao) DeleteMenu(ctx context.Context, menuID uint) error {
	return menuDao.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&system.SysMenu{}, "parent_id = ?", menuID).Error; err == nil {
			// 存在子菜单，不能删除
			return errors.WithCode(errcode.ErrMenuHasChildren, "")
		}

		var menu system.SysMenu
		if err := tx.First(&menu, menuID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.WithCode(errcode.ErrMenuNotFound, "%s", err.Error())
			}
			return errors.WithCode(errcode.ErrDatabase, "%s", err.Error())
		}

		if err := tx.First(&system.SysRoleAuthority{}, "default_router = ?", menu.Name).Error; err == nil {
			// 被角色作为首页，不能删除
			return errors.WithCode(errcode.ErrMenuAsHomePage, "")
		}

		if err := tx.Delete(&system.SysMenu{}, "id = ?", menuID).Error; err != nil {
			return errors.WithCode(errcode.ErrMenuDeleteFailed, "%s", err.Error())
		}

		// 删除权限菜单关联
		if err := tx.Delete(&system.SysMenuAuthorityRelation{}, "sys_base_menu_id = ?", menuID).Error; err != nil {
			return errors.WithCode(errcode.ErrMenuDeleteFailed, "%s", err.Error())
		}

		return nil
	})

}

// getMenuTree 根据权限ID构建菜单树结构
// 参数：
//
//	authorityID: 用户权限ID
//
// 返回值：
//
//	treeMap: 菜单树映射，key为父菜单ID，value为该父菜单下的所有子菜单列表
//	err: 错误信息
//
// 功能：
//
//	通过查询数据库获取用户权限对应的菜单信息，并组织成树形结构便于前端展示
func (menuDao *MenuDao) getMenuTree(ctx context.Context, authorityID uint) (treeMap map[uint][]system.SysAuthorizedMenu, err error) {
	// 初始化变量
	var allMenus []system.SysAuthorizedMenu // 存储最终生成的菜单列表
	var baseMenu []system.SysMenu           // 存储从数据库查询的基础菜单
	// var btns []system.SysAuthorityBtn         // 存储按钮权限(暂未使用)

	// 1. 查询该角色ID对应的所有菜单权限关联记录
	// 通过SysMenuAuthorityRelation表查询指定角色拥有的所有菜单权限
	var SysAuthorityMenus []system.SysMenuAuthorityRelation

	if err = menuDao.GetDB(ctx).Where("authority_id = ?", authorityID).Find(&SysAuthorityMenus).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(errcode.ErrAuthorityMenuNotFound, "%s", err.Error())
		}
		return nil, errors.WithCode(errcode.ErrDatabase, "%s", err.Error())
	}

	// 2. 提取所有菜单ID
	// 从权限菜单关联记录中提取出所有菜单的ID，用于后续查询具体菜单信息
	var MenuIDs []string
	for i := range SysAuthorityMenus {
		MenuIDs = append(MenuIDs, SysAuthorityMenus[i].MenuID)
	}

	// 3. 根据菜单ID查询基础菜单信息（按sort排序并预加载参数）
	// 查询具体的菜单信息，按照sort字段排序以便正确显示菜单顺序
	err = menuDao.GetDB(ctx).Where("id in (?)", MenuIDs).Order("sort").Find(&baseMenu).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(errcode.ErrMenuNotFound, "%s", err.Error())
		}
		return nil, errors.WithCode(errcode.ErrDatabase, "%s", err.Error())
	}

	// 4. 将基础菜单转换为系统菜单格式
	// 将查询到的基础菜单信息转换为SysMenu结构，添加角色ID等额外信息
	for i := range baseMenu {
		allMenus = append(allMenus, system.SysAuthorizedMenu{
			SysMenu:     baseMenu[i],    // 基础菜单信息
			AuthorityID: authorityID,    // 所属角色ID
			MenuID:      baseMenu[i].ID, // 菜单ID
		})
	}

	// 5. 将按钮权限分配到对应菜单，并构建菜单树结构
	// 遍历所有菜单，按照父菜单ID进行分组，构建菜单树映射结构
	// treeMap的key是父菜单ID，value是具有相同父菜单ID的所有子菜单列表
	treeMap = make(map[uint][]system.SysAuthorizedMenu)
	for _, v := range allMenus {
		treeMap[v.ParentID] = append(treeMap[v.ParentID], v) // 按父ID分组，构建树形结构
	}

	// 返回构建好的菜单树映射和nil错误表示操作成功
	return treeMap, nil
}

// getChildrenList 递归获取菜单的子菜单列表
// 参数：
//
//	menu: 指向当前菜单项的指针，函数将为其填充子菜单信息
//	treeMap: 菜单树映射，key为父菜单ID，value为该父菜单下的所有子菜单列表
//
// 返回值：
//
//	err: 错误信息
//
// 功能：
//
//	通过递归方式为指定菜单项填充其所有的子菜单信息，构建完整的菜单树结构
func (menuDao *MenuDao) getChildrenList(menu *system.SysAuthorizedMenu, treeMap map[uint][]system.SysAuthorizedMenu) (err error) {
	// 从treeMap中获取当前菜单项的所有直接子菜单，并赋值给menu.Children
	// treeMap的key是父菜单ID（即menu.MenuID），value是具有相同父菜单ID的所有子菜单列表
	menu.Children = treeMap[menu.MenuID]

	// 遍历当前菜单的所有直接子菜单
	for i := range menu.Children {
		// 对每个子菜单递归调用getChildrenList方法，继续获取其子菜单列表
		// 这样可以构建出完整的多级菜单树结构
		err = menuDao.getChildrenList(&menu.Children[i], treeMap)

		// 如果在递归过程中发生错误，则直接返回错误
		if err != nil {
			return err
		}
	}

	// 返回nil错误表示操作成功完成
	return err
}

// GetMenuList 根据权限ID获取菜单信息列表
// 参数：
//
//	authorityID: 用户权限ID
//
// 返回值：
//
//	menuList: 菜单信息列表
//	err: 错误信息
//
// 功能：
//
//	通过构建完整的菜单树结构来返回用户有权访问的菜单信息
func (menuDao *MenuDao) GetMenuList(ctx context.Context, parentAuthorityID, authorityID uint) (menuList []system.SysMenu, err error) {
	// 调用getBaseMenuTreeMap方法获取以权限ID为基础的菜单树映射结构
	// treeMap的key是父菜单ID，value是具有相同父菜单ID的所有子菜单列表
	treeMap, err := menuDao.getBaseMenuTreeMap(ctx, parentAuthorityID, authorityID)

	// 如果获取菜单树映射失败，直接返回错误
	if err != nil {

		return nil, err
	}

	// 从菜单树映射中获取顶级菜单项列表（父ID为0的菜单项）
	// menuList包含所有顶级菜单项
	menuList = treeMap[0]

	// 遍历所有顶级菜单项，为每个菜单项递归填充其子菜单列表
	for i := 0; i < len(menuList); i++ {
		// 调用getBaseChildrenList方法递归获取并设置当前菜单项的完整子菜单树
		err = menuDao.getBaseChildrenList(&menuList[i], treeMap)

		// 如果在获取子菜单过程中出现错误，则中断循环并返回错误
		if err != nil {
			return nil, err
		}
	}

	// 返回完整的菜单列表和nil错误表示操作成功
	return menuList, nil
}

// getBaseMenuTreeMap 根据权限ID获取基础菜单树映射
// 参数：
//
//	authorityID: 用户权限ID
//
// 返回值：
//
//	treeMap: 菜单树映射，key为父菜单ID，value为该父菜单下的所有子菜单列表
//	err: 错误信息
//
// 功能：
//
//	通过查询数据库获取菜单信息，并根据权限配置决定是否进行菜单筛选，最终组织成树形结构映射
func (menuDao *MenuDao) getBaseMenuTreeMap(ctx context.Context, parentAuthorityID, authorityID uint) (treeMap map[uint][]system.SysMenu, err error) {
	if err != nil {
		// 如果获取父权限ID失败，返回错误
		return nil, err
	}

	// 初始化变量
	var allMenus []system.SysMenu // 存储查询到的所有基础菜单
	// 创建数据库查询对象，按sort字段排序，并预加载菜单按钮和参数信息
	// db := global.MPA_DB.Order("sort").Preload("MenuBtn").Preload("Parameters")
	db := menuDao.GetDB(ctx).Order("sort")

	// 当开启了严格的树角色并且父角色不为0时需要进行菜单筛选
	// 严格的树角色意味着子角色只能访问父角色拥有的菜单的子集
	if global.MPA_CONFIG.System.UseStrictAuth && parentAuthorityID != 0 {
		// 查询当前权限ID对应的所有菜单权限关联记录
		var authorityMenus []system.SysMenuAuthorityRelation
		err = menuDao.GetDB(ctx).Where("sys_authority_authority_id = ?", authorityID).Find(&authorityMenus).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.WithCode(errcode.ErrAuthorityMenuNotFound, "%s", err.Error())
			}
			return nil, errors.WithCode(errcode.ErrDatabase, "%s", err.Error())
		}

		// 提取所有菜单ID
		var menuIDs []string
		for i := range authorityMenus {
			menuIDs = append(menuIDs, authorityMenus[i].MenuID)
		}

		// 在基础查询条件上增加菜单ID筛选条件
		db = db.Where("id in (?)", menuIDs)
	}

	// 执行数据库查询，获取所有符合条件的基础菜单
	err = db.Find(&allMenus).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(errcode.ErrMenuNotFound, "%s", err.Error())
		}
		return nil, errors.WithCode(errcode.ErrDatabase, "%s", err.Error())
	}
	// 构建菜单树映射结构
	// 遍历所有菜单，按照父菜单ID进行分组
	// treeMap的key是父菜单ID，value是具有相同父菜单ID的所有子菜单列表
	treeMap = make(map[uint][]system.SysMenu)
	for _, v := range allMenus {
		treeMap[v.ParentID] = append(treeMap[v.ParentID], v) // 按父ID分组，构建树形结构映射
	}

	// 返回构建好的菜单树映射和nil错误表示操作成功
	return treeMap, err
}

// getBaseChildrenList 递归获取基础菜单的子菜单列表
// 参数：
//
//	menu: 指向当前菜单项的指针，函数将为其填充子菜单信息
//	treeMap: 菜单树映射，key为父菜单ID，value为该父菜单下的所有子菜单列表
//
// 返回值：
//
//	err: 错误信息
//
// 功能：
//
//	通过递归方式为指定菜单项填充其所有的子菜单信息，构建完整的菜单树结构
func (menuDao *MenuDao) getBaseChildrenList(menu *system.SysMenu, treeMap map[uint][]system.SysMenu) (err error) {
	// 从treeMap中获取当前菜单项的所有直接子菜单，并赋值给menu.Children
	// treeMap的key是父菜单ID，value是具有相同父菜单ID的所有子菜单列表
	menu.Children = treeMap[menu.ID]

	// 遍历当前菜单的所有直接子菜单
	for i := 0; i < len(menu.Children); i++ {
		// 对每个子菜单递归调用getBaseChildrenList方法，继续获取其子菜单列表
		// 这样可以构建出完整的多级菜单树结构
		err = menuDao.getBaseChildrenList(&menu.Children[i], treeMap)

		// 如果在递归过程中发生错误，则直接返回错误
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.WithCode(errcode.ErrMenuNotFound, "%s", err.Error())
			}
			return errors.WithCode(errcode.ErrDatabase, "%s", err.Error())
		}
	}

	// 返回nil错误表示操作成功完成
	return nil
}

// GetMenuAuthority 根据权限ID获取菜单权限列表
// 参数：
//
//	authorityID: 权限ID
//
// 返回值：
//
//	menus: 权限对应的菜单列表
//	err: 错误信息
//
// 功能：
//
//	获取指定权限ID对应的所有菜单权限信息，通常用于权限管理中展示某个角色拥有的菜单权限
func (menuDao *MenuDao) GetMenuAuthority(ctx context.Context, authorityID uint) (menus []system.SysAuthorizedMenu, err error) {
	// 声明基础菜单切片，用于存储从数据库查询到的基础菜单信息
	var baseMenu []system.SysMenu

	// 声明权限菜单关联切片，用于存储权限与菜单的关联关系
	var SysAuthorityMenus []system.SysMenuAuthorityRelation

	// 查询sys_authority_menu表，获取指定权限ID关联的所有菜单记录
	// Where条件: sys_authority_authority_id等于传入的权限ID
	// Find方法将查询结果填充到SysAuthorityMenus切片中
	// Error属性用于获取查询过程中的错误信息
	err = menuDao.GetDB(ctx).Where("sys_authority_authority_id = ?", authorityID).Find(&SysAuthorityMenus).Error

	// 检查查询过程中是否出现错误
	if err != nil {
		// 如果查询出错，直接返回错误信息
		return
	}

	// 声明菜单ID切片，用于存储从权限菜单关联记录中提取的菜单ID
	var MenuIDs []string

	// 遍历权限菜单关联记录，提取所有菜单ID
	for i := range SysAuthorityMenus {
		// 将每个关联记录的MenuID添加到MenuIDs切片中
		MenuIDs = append(MenuIDs, SysAuthorityMenus[i].MenuID)
	}

	// 根据菜单ID列表查询sys_base_menu表获取具体的菜单信息
	// Where条件: id在MenuIDs切片中的所有值
	// Order条件: 按照sort字段升序排列
	// Find方法将查询结果填充到baseMenu切片中
	// 注意：如果MenuIDs为空，此查询可能不会返回任何结果
	err = menuDao.GetDB(ctx).Where("id in (?) ", MenuIDs).Order("sort").Find(&baseMenu).Error

	// 遍历基础菜单列表，将其转换为系统菜单格式
	for i := range baseMenu {
		// 将每个基础菜单转换为系统菜单并添加到menus切片中
		menus = append(menus, system.SysAuthorizedMenu{
			// 继承基础菜单的所有字段
			SysMenu: baseMenu[i],
			// 设置权限ID为传入的权限ID
			AuthorityID: authorityID,
			// 设置菜单ID为基础菜单的ID
			MenuID: baseMenu[i].ID,
			// 注释掉的Parameters字段可能是历史遗留代码或待实现功能
			// Parameters:  baseMenu[i].Parameters,
		})
	}

	// 返回转换后的菜单列表和可能的错误信息
	// 如果前面的查询都成功，err应该为nil
	return menus, err
}
