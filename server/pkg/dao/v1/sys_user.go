// Package v1 提供数据访问层的实现
// 包含用户、API、权限、菜单、Casbin等数据的CRUD操作
package v1

import (
	"context"
	"time"

	"github.com/marmotedu/errors"

	"server/pkg/api/request"
	"server/pkg/errcode"
	"server/pkg/model/system"
	"server/pkg/utils"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Register 注册新用户
// 参数：
//
//	u: 用户信息
//
// 返回值：
//
//	userInter: 注册成功的用户信息
//	err: 错误信息
func (userDao *UserDao) Register(ctx context.Context, u system.SysUser) (userInter system.SysUser, err error) {
	var user system.SysUser
	if !errors.Is(userDao.GetDB(ctx).Where("user_name = ?", u.UserName).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return userInter, errors.WithCode(errcode.ErrUserAlreadyExist, "用户名已注册")
	}
	// 否则 附加uuid 密码hash加密 注册
	u.Password = utils.BcryptHash(u.Password)
	u.UUID = uuid.NewV4()
	if err = userDao.GetDB(ctx).Create(&u).Error; err != nil {
		return userInter, errors.WithCode(errcode.ErrUserRegisteFailed, "%s", err.Error())
	}
	return u, nil
}

// Login 用户登录
// 参数：
//
//	u: 包含用户名和密码的用户信息
//
// 返回值：
//
//	userInter: 登录成功的用户信息
//	err: 错误信息
func (userDao *UserDao) Login(ctx context.Context, u system.SysUser) (userInter *system.SysUser, err error) {
	var user system.SysUser
	err = userDao.GetDB(ctx).Where("user_name = ?", u.UserName).First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return nil, errors.WithCode(errcode.ErrUserPasswordFault, "密码错误")
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.WithCode(errcode.ErrUserNotFound, "用户不存在")
	}

	return &user, err
}

// ChangePassword 修改用户密码
// 参数：
//
//	u: 包含用户ID和旧密码的用户信息
//	newPassword: 新密码
//
// 返回值：
//
//	userInter: 修改密码成功的用户信息
//	err: 错误信息
func (userDao *UserDao) ChangePassword(ctx context.Context, u system.SysUser, newPassword string) (userInter *system.SysUser, err error) {
	var user system.SysUser
	if err = userDao.GetDB(ctx).Where("id = ?", u.ID).First(&user).Error; err != nil {
		return nil, err
	}
	if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
		return nil, errors.WithCode(errcode.ErrUserPasswordFault, "密码错误")
	}
	user.Password = utils.BcryptHash(newPassword)

	if err = userDao.GetDB(ctx).Save(&user).Error; err != nil {
		return nil, errors.WithCode(errcode.ErrPasswordIncorrect, "%s", err.Error())
	}

	return &user, nil
}

// GetUserInfoList 分页获取用户列表
// 参数：
//
//	info: 分页和查询条件信息
//	authorityId: 权限ID
//
// 返回值：
//
//	list: 用户列表
//	total: 总记录数
//	err: 错误信息
func (userDao *UserDao) GetUserInfoList(ctx context.Context, info request.GetUserList, authorityID uint) (list []system.SysUser, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	var authorityIDs []uint
	// 查看当前角色及以下角色的角色ID
	if err = userDao.GetDB(ctx).Model(&system.SysRoleAuthority{}).Select("authority_id").Where("parent_id = ?", authorityID).Find(&authorityIDs).Error; err != nil {
		return list, total, errors.WithCode(errcode.ErrAuthorityNotPermissionGet, "%s", err.Error())
	}

	var userIDs []uint
	// 查看当前角色及以下角色的用户ID
	if err = userDao.GetDB(ctx).Model(&system.SysUserAuthorityRelation{}).Select("user_id").Where("authority_id IN ?", authorityIDs).Find(&userIDs).Error; err != nil {
		return list, total, errors.WithCode(errcode.ErrAuthorityNotPermissionGet, "%s", err.Error())
	}

	db := userDao.GetDB(ctx).Model(&system.SysUser{}).Preload("Authorities").Where("id IN ?", userIDs)

	// 用户名模糊查询
	if info.Username != "" {
		db = db.Where("user_name LIKE ?", "%"+info.Username+"%")
	}
	// 昵称模糊查询
	if info.NickName != "" {
		db = db.Where("nick_name LIKE ?", "%"+info.NickName+"%")
	}
	// 手机号模糊查询
	if info.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+info.Phone+"%")
	}
	// 邮箱模糊查询
	if info.Email != "" {
		db = db.Where("email LIKE ?", "%"+info.Email+"%")
	}

	var userList []system.SysUser

	if err = db.Count(&total).Error; err != nil {
		return
	}

	if err = db.Preload("Authority").Limit(limit).Offset(offset).Find(&userList).Error; err != nil {
		return list, total, errors.WithCode(errcode.ErrUserSetInfoFailed, "%s", err.Error())
	}

	return userList, total, err
}

// SetUserAuthority 设置用户权限
// 参数：
//
//	id: 用户ID
//	authorityID: 权限ID
//
// 返回值：
//
//	err: 错误信息
func (userDao *UserDao) SetUserAuthority(ctx context.Context, id uint, authorityID uint) (err error) {
	assignErr := userDao.GetDB(ctx).Where("user_id = ? AND authority_id = ?", id, authorityID).First(&system.SysUserAuthorityRelation{}).Error
	if errors.Is(assignErr, gorm.ErrRecordNotFound) {
		return errors.New("该用户无此角色")
	}
	err = userDao.GetDB(ctx).Where("id = ?", id).First(&system.SysUser{}).Update("authority_id", authorityID).Error
	return err
}

// SetUserAuthorities 设置用户多个权限
// 参数：
//
//	id: 用户ID
//	authorityIDs: 权限ID列表
//
// 返回值：
//
//	err: 错误信息
func (userDao *UserDao) SetUserAuthorities(ctx context.Context, id uint, authorityIDs []uint) (err error) {
	return userDao.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		if txErr := tx.Delete(&[]system.SysUserAuthorityRelation{}, "user_id = ?", id).Error; txErr != nil {
			return errors.WithCode(errcode.ErrUserSetAuthorities, "%s", txErr.Error())
		}

		var useAuthority []system.SysUserAuthorityRelation
		for _, v := range authorityIDs {
			useAuthority = append(useAuthority, system.SysUserAuthorityRelation{
				UserID:      id,
				AuthorityID: v,
			})
		}

		if txErr := tx.Create(&useAuthority).Error; txErr != nil {
			return errors.WithCode(errcode.ErrUserSetAuthorities, "%s", txErr.Error())
		}

		if txErr := tx.Where("id = ?", id).First(&system.SysUser{}).Update("authority_id", authorityIDs[0]).Error; txErr != nil {
			return errors.WithCode(errcode.ErrUserSetAuthorities, "%s", txErr.Error())
		}

		// 返回 nil 提交事务
		return nil
	})
}

// DeleteUser 删除用户
// 参数：
//
//	id: 用户ID
//
// 返回值：
//
//	err: 错误信息
func (userDao *UserDao) DeleteUser(ctx context.Context, id int) (err error) {
	var user system.SysUser

	return userDao.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		if txErr := tx.Where("id = ?", id).Delete(&user).Error; txErr != nil {
			return errors.WithCode(errcode.ErrUserDeleteFailed, "%s", txErr.Error())
		}

		if txErr := tx.Delete(&[]system.SysUserAuthorityRelation{}, "user_id = ?", id).Error; txErr != nil {
			return errors.WithCode(errcode.ErrUserDeleteFailed, "%s", txErr.Error())
		}

		return nil
	})
}

// SetUserInfo 设置用户信息
// 参数：
//
//	req: 用户信息
//
// 返回值：
//
//	err: 错误信息
func (userDao *UserDao) SetUserInfo(ctx context.Context, req system.SysUser) error {
	if err := userDao.GetDB(ctx).Model(&system.SysUser{}).
		Select("updated_at", "nick_name", "avatar", "phone", "email", "is_active").
		Where("id=?", req.ID).
		Updates(map[string]any{
			"updated_at": time.Now(),
			"nick_name":  req.NickName,
			"avatar":     req.Avatar,
			"phone":      req.Phone,
			"email":      req.Email,
			"is_active":  req.IsActive,
		}).Error; err != nil {
		return errors.WithCode(errcode.ErrUserSetInfoFailed, "%s", err.Error())
	}

	return nil
}

// SetSelfInfo 设置用户自身信息
// 参数：
//
//	req: 用户信息
//
// 返回值：
//
//	err: 错误信息
func (userDao *UserDao) SetSelfInfo(ctx context.Context, req system.SysUser) error {
	if err := userDao.GetDB(ctx).Model(&system.SysUser{}).
		Where("id=?", req.ID).
		Updates(req).Error; err != nil {
		return errors.WithCode(errcode.ErrUserSetSelfFailed, "%s", err.Error())
	}

	return nil
}

// GetUserInfo 根据UUID获取用户信息
// 参数：
//
//	uuid: 用户UUID
//
// 返回值：
//
//	reqUser: 用户信息
//	err: 错误信息
func (userDao *UserDao) GetUserInfo(ctx context.Context, uuid uuid.UUID) (reqUser system.SysUser, err error) {
	if err = userDao.GetDB(ctx).Preload("Authorities").Preload("Authority").First(&reqUser, "uuid = ?", uuid).Error; err != nil {
		return reqUser, errors.WithCode(errcode.ErrUserNotFound, "%s", errors.New("用户不存在"))
	}

	return reqUser, nil
}

// FindUserByID 通过ID获取用户信息
// 参数：
//
//	id: 用户ID
//
// 返回值：
//
//	user: 用户信息
//	err: 错误信息
func (userDao *UserDao) FindUserByID(ctx context.Context, id int) (user *system.SysUser, err error) {
	var u system.SysUser
	err = userDao.GetDB(ctx).Where("`id` = ?", id).First(&u).Error
	return &u, err
}

// FindUserByUUID 通过UUID获取用户信息
// 参数：
//
//	uuid: 用户UUID
//
// 返回值：
//
//	user: 用户信息
//	err: 错误信息
func (userDao *UserDao) FindUserByUUID(ctx context.Context, uuid string) (user *system.SysUser, err error) {
	var u system.SysUser
	if err = userDao.GetDB(ctx).Where("`uuid` = ?", uuid).First(&u).Error; err != nil {
		return &u, errors.New("用户不存在")
	}
	return &u, nil
}

// ResetPassword 重置用户密码
// 参数：
//
//	ID: 用户ID
//
// 返回值：
//
//	err: 错误信息
func (userDao *UserDao) ResetPassword(ctx context.Context, ID uint) error {
	if err := userDao.GetDB(ctx).Model(&system.SysUser{}).Where("id = ?", ID).Update("password", utils.BcryptHash("123456")).Error; err != nil {
		return errors.WithCode(errcode.ErrUserResetPassworkFailed, "%s", err.Error())
	}
	return nil
}

// UpdateLoginTime 更新用户的登录时间
// 参数：
//
//	id: 用户ID
//
// 返回值：
//
//	err: 错误信息
func (userDao *UserDao) UpdateLoginTime(ctx context.Context, id uint) error {
	if err := userDao.GetDB(ctx).Model(&system.SysUser{}).Where("id = ?", id).Update("last_login", time.Now().Local()).Error; err != nil {
		return errors.WithCode(errcode.ErrUserSetInfoFailed, "%s", err.Error())
	}
	return nil
}
