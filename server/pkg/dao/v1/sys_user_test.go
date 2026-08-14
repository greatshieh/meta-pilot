package v1

import (
	"context"
	"testing"
	"time"

	"github.com/marmotedu/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"server/pkg/api/request"
	"server/pkg/errcode"
	common_req "server/pkg/model/common/request"
	"server/pkg/model/system"
)

func setupTestContainer(t *testing.T) (context.Context, *gorm.DB, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "password",
			"MYSQL_DATABASE":      "test_db",
		},
		WaitingFor: wait.ForLog("port: 3306  MySQL Community Server - GPL"),
	}

	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)

	ip, err := mysqlC.Host(ctx)
	assert.NoError(t, err)

	port, err := mysqlC.MappedPort(ctx, "3306")
	assert.NoError(t, err)

	dsn := "root:password@tcp(" + ip + ":" + port.Port() + ")/test_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	assert.NoError(t, err)

	err = db.AutoMigrate(
		&system.SysUser{},
		&system.SysRoleAuthority{},
		&system.SysUserAuthorityRelation{},
	)
	assert.NoError(t, err)

	cleanup := func() {
		mysqlC.Terminate(ctx)
	}

	return ctx, db, cleanup
}

func TestUserDao_Register(t *testing.T) {
	ctx, db, cleanup := setupTestContainer(t)
	defer cleanup()

	userDao := NewUserDao(db)

	t.Run("register new user successfully", func(t *testing.T) {
		user := system.SysUser{
			UserName: "testuser_new",
			Password: "password123",
			NickName: "Test User",
		}

		registeredUser, err := userDao.Register(ctx, user)

		assert.NoError(t, err)
		assert.NotEmpty(t, registeredUser.ID)
		assert.Equal(t, "testuser_new", registeredUser.UserName)
		assert.NotEqual(t, "password123", registeredUser.Password)
		assert.NotEmpty(t, registeredUser.UUID)
		assert.True(t, registeredUser.IsActive)
	})

	t.Run("register existing user", func(t *testing.T) {
		user := system.SysUser{
			UserName: "testuser_new",
			Password: "password456",
			NickName: "Test User 2",
		}

		registeredUser, err := userDao.Register(ctx, user)
		assert.Error(t, err)
		assert.Empty(t, registeredUser.ID)
		assert.True(t, errors.IsCode(err, errcode.ErrUserAlreadyExist))
	})

	t.Run("register user with empty username", func(t *testing.T) {
		user := system.SysUser{
			UserName: "",
			Password: "password123",
			NickName: "Empty Username",
		}

		registeredUser, err := userDao.Register(ctx, user)

		assert.NoError(t, err)
		assert.NotEmpty(t, registeredUser.ID)
		assert.Equal(t, "", registeredUser.UserName)
	})

	t.Run("register user with empty password", func(t *testing.T) {
		user := system.SysUser{
			UserName: "testuser_empty_pwd",
			Password: "",
			NickName: "Empty Password",
		}

		registeredUser, err := userDao.Register(ctx, user)

		assert.NoError(t, err)
		assert.NotEmpty(t, registeredUser.ID)
		assert.NotEmpty(t, registeredUser.Password)
	})

	t.Run("register user with context timeout", func(t *testing.T) {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		time.Sleep(10 * time.Millisecond)

		user := system.SysUser{
			UserName: "testuser_timeout",
			Password: "password123",
			NickName: "Timeout Test",
		}

		_, err := userDao.Register(timeoutCtx, user)

		assert.Error(t, err)
		assert.True(t, errors.Is(timeoutCtx.Err(), context.DeadlineExceeded))
	})
}

func TestUserDao_Login(t *testing.T) {
	ctx, db, cleanup := setupTestContainer(t)
	defer cleanup()

	userDao := NewUserDao(db)

	registerUser := system.SysUser{
		UserName: "login_test_user",
		Password: "login_password",
		NickName: "Login Test",
	}
	_, err := userDao.Register(ctx, registerUser)
	assert.NoError(t, err)

	t.Run("login with correct password", func(t *testing.T) {
		user := system.SysUser{
			UserName: "login_test_user",
			Password: "login_password",
		}

		loggedInUser, err := userDao.Login(ctx, user)

		assert.NoError(t, err)
		assert.NotNil(t, loggedInUser)
		assert.Equal(t, "login_test_user", loggedInUser.UserName)
	})

	t.Run("login with wrong password", func(t *testing.T) {
		user := system.SysUser{
			UserName: "login_test_user",
			Password: "wrong_password",
		}

		loggedInUser, err := userDao.Login(ctx, user)

		assert.Error(t, err)
		assert.Nil(t, loggedInUser)
		assert.True(t, errors.IsCode(err, errcode.ErrUserPasswordFault))
	})

	t.Run("login with non-existent user", func(t *testing.T) {
		user := system.SysUser{
			UserName: "non_existent_user",
			Password: "any_password",
		}

		loggedInUser, err := userDao.Login(ctx, user)

		assert.Error(t, err)
		assert.Nil(t, loggedInUser)
		assert.True(t, errors.IsCode(err, errcode.ErrUserNotFound))
	})

	t.Run("login with empty username", func(t *testing.T) {
		user := system.SysUser{
			UserName: "",
			Password: "any_password",
		}

		loggedInUser, err := userDao.Login(ctx, user)

		assert.Error(t, err)
		assert.Nil(t, loggedInUser)
		assert.True(t, errors.IsCode(err, errcode.ErrUserNotFound))
	})
}

func TestUserDao_ChangePassword(t *testing.T) {
	ctx, db, cleanup := setupTestContainer(t)
	defer cleanup()

	userDao := NewUserDao(db)

	registerUser := system.SysUser{
		UserName: "change_pwd_user",
		Password: "old_password",
		NickName: "Change Password Test",
	}
	registered, err := userDao.Register(ctx, registerUser)
	assert.NoError(t, err)

	t.Run("change password successfully", func(t *testing.T) {
		user := system.SysUser{}
		user.ID = registered.ID
		user.Password = "old_password"

		updatedUser, err := userDao.ChangePassword(ctx, user, "new_password")

		assert.NoError(t, err)
		assert.NotNil(t, updatedUser)
		assert.Equal(t, registered.ID, updatedUser.ID)
	})

	t.Run("change password with wrong old password", func(t *testing.T) {
		user := system.SysUser{}
		user.ID = registered.ID
		user.Password = "wrong_old_password"

		updatedUser, err := userDao.ChangePassword(ctx, user, "new_password")

		assert.Error(t, err)
		assert.Nil(t, updatedUser)
		assert.True(t, errors.IsCode(err, errcode.ErrUserPasswordFault))
	})

	t.Run("change password for non-existent user", func(t *testing.T) {
		user := system.SysUser{}
		user.ID = 99999
		user.Password = "any_password"

		updatedUser, err := userDao.ChangePassword(ctx, user, "new_password")

		assert.Error(t, err)
		assert.Nil(t, updatedUser)
	})
}

func TestUserDao_GetUserInfoList(t *testing.T) {
	ctx, db, cleanup := setupTestContainer(t)
	defer cleanup()

	userDao := NewUserDao(db)

	authority := system.SysRoleAuthority{
		AuthorityID:   1,
		AuthorityName: "Test Authority",
	}
	err := db.Create(&authority).Error
	assert.NoError(t, err)

	roleAuthority := system.SysRoleAuthority{
		AuthorityID:   2,
		AuthorityName: "Sub Authority",
		ParentID:      &authority.AuthorityID,
	}
	err = db.Create(&roleAuthority).Error
	assert.NoError(t, err)

	for i := 0; i < 5; i++ {
		user := system.SysUser{
			UserName:    "list_user_" + string(rune('a'+i)),
			Password:    "password123",
			NickName:    "List User " + string(rune('a'+i)),
			AuthorityID: roleAuthority.AuthorityID,
		}
		registered, err := userDao.Register(ctx, user)
		assert.NoError(t, err)

		userAuth := system.SysUserAuthorityRelation{
			UserID:      registered.ID,
			AuthorityID: roleAuthority.AuthorityID,
		}
		err = db.Create(&userAuth).Error
		assert.NoError(t, err)
	}

	t.Run("get user list with pagination", func(t *testing.T) {
		info := request.GetUserList{
			PageInfo: common_req.PageInfo{
				Page:     1,
				PageSize: 2,
			},
		}

		list, _, err := userDao.GetUserInfoList(ctx, info, 1)

		assert.NoError(t, err)
		assert.Len(t, list, 2)
	})

	t.Run("get user list with username filter", func(t *testing.T) {
		info := request.GetUserList{
			PageInfo: common_req.PageInfo{
				Page:     1,
				PageSize: 10,
			},
			Username: "list_user_a",
		}

		list, total, err := userDao.GetUserInfoList(ctx, info, 1)

		assert.NoError(t, err)
		assert.Len(t, list, 1)
		assert.Equal(t, int64(1), total)
	})

	t.Run("get user list with nickname filter", func(t *testing.T) {
		info := request.GetUserList{
			PageInfo: common_req.PageInfo{
				Page:     1,
				PageSize: 10,
			},
			NickName: "List User",
		}

		_, _, err := userDao.GetUserInfoList(ctx, info, 1)

		assert.NoError(t, err)
	})
}

func TestUserDao_SetUserAuthority(t *testing.T) {
	ctx, db, cleanup := setupTestContainer(t)
	defer cleanup()

	userDao := NewUserDao(db)

	user := system.SysUser{
		UserName: "authority_user",
		Password: "password123",
		NickName: "Authority Test",
	}
	registered, err := userDao.Register(ctx, user)
	assert.NoError(t, err)

	userAuth := system.SysUserAuthorityRelation{
		UserID:      registered.ID,
		AuthorityID: 1,
	}
	err = db.Create(&userAuth).Error
	assert.NoError(t, err)

	t.Run("set user authority successfully", func(t *testing.T) {
		err := userDao.SetUserAuthority(ctx, registered.ID, 1)

		assert.NoError(t, err)
	})

	t.Run("set user authority without permission", func(t *testing.T) {
		err := userDao.SetUserAuthority(ctx, registered.ID, 999)

		assert.Error(t, err)
		assert.Equal(t, "该用户无此角色", err.Error())
	})
}

func TestUserDao_SetUserAuthorities(t *testing.T) {
	ctx, db, cleanup := setupTestContainer(t)
	defer cleanup()

	userDao := NewUserDao(db)

	user := system.SysUser{
		UserName: "multi_authority_user",
		Password: "password123",
		NickName: "Multi Authority Test",
	}
	registered, err := userDao.Register(ctx, user)
	assert.NoError(t, err)

	t.Run("set user authorities successfully", func(t *testing.T) {
		err := userDao.SetUserAuthorities(ctx, registered.ID, []uint{1, 2, 3})

		assert.NoError(t, err)
	})
}

func TestUserDao_DeleteUser(t *testing.T) {
	ctx, db, cleanup := setupTestContainer(t)
	defer cleanup()

	userDao := NewUserDao(db)

	user := system.SysUser{
		UserName: "delete_user",
		Password: "password123",
		NickName: "Delete Test",
	}
	registered, err := userDao.Register(ctx, user)
	assert.NoError(t, err)

	t.Run("delete user successfully", func(t *testing.T) {
		err := userDao.DeleteUser(ctx, int(registered.ID))

		assert.NoError(t, err)
	})

	t.Run("delete non-existent user", func(t *testing.T) {
		err := userDao.DeleteUser(ctx, 99999)

		assert.NoError(t, err)
	})
}

func TestUserDao_SetUserInfo(t *testing.T) {
	ctx, db, cleanup := setupTestContainer(t)
	defer cleanup()

	userDao := NewUserDao(db)

	user := system.SysUser{
		UserName: "setinfo_user",
		Password: "password123",
		NickName: "Set Info Test",
	}
	registered, err := userDao.Register(ctx, user)
	assert.NoError(t, err)

	t.Run("set user info successfully", func(t *testing.T) {
		updateUser := system.SysUser{}
		updateUser.ID = registered.ID
		updateUser.NickName = "Updated Nickname"
		updateUser.Phone = "12345678900"
		updateUser.Email = "test@example.com"
		updateUser.IsActive = false

		err := userDao.SetUserInfo(ctx, updateUser)

		assert.NoError(t, err)

		var updated system.SysUser
		err = db.First(&updated, registered.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, "Updated Nickname", updated.NickName)
		assert.Equal(t, "12345678900", updated.Phone)
		assert.False(t, updated.IsActive)
	})
}

func TestUserDao_SetSelfInfo(t *testing.T) {
	ctx, db, cleanup := setupTestContainer(t)
	defer cleanup()

	userDao := NewUserDao(db)

	user := system.SysUser{
		UserName: "setself_user",
		Password: "password123",
		NickName: "Set Self Info Test",
	}
	registered, err := userDao.Register(ctx, user)
	assert.NoError(t, err)

	t.Run("set self info successfully", func(t *testing.T) {
		updateUser := system.SysUser{}
		updateUser.ID = registered.ID
		updateUser.NickName = "Self Updated"

		err := userDao.SetSelfInfo(ctx, updateUser)

		assert.NoError(t, err)
	})
}

func TestUserDao_GetUserInfo(t *testing.T) {
	ctx, db, cleanup := setupTestContainer(t)
	defer cleanup()

	userDao := NewUserDao(db)

	user := system.SysUser{
		UserName: "getinfo_user",
		Password: "password123",
		NickName: "Get Info Test",
	}
	registered, err := userDao.Register(ctx, user)
	assert.NoError(t, err)

	t.Run("get user info by uuid", func(t *testing.T) {
		foundUser, err := userDao.GetUserInfo(ctx, registered.UUID)

		assert.NoError(t, err)
		assert.Equal(t, registered.UserName, foundUser.UserName)
	})

	t.Run("get user info with non-existent uuid", func(t *testing.T) {
		nonExistentUUID, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")
		_, err := userDao.GetUserInfo(ctx, nonExistentUUID)

		assert.Error(t, err)
		assert.True(t, errors.IsCode(err, errcode.ErrUserNotFound))
	})
}

func TestUserDao_FindUserByID(t *testing.T) {
	ctx, db, cleanup := setupTestContainer(t)
	defer cleanup()

	userDao := NewUserDao(db)

	user := system.SysUser{
		UserName: "findbyid_user",
		Password: "password123",
		NickName: "Find By ID Test",
	}
	registered, err := userDao.Register(ctx, user)
	assert.NoError(t, err)

	t.Run("find user by id", func(t *testing.T) {
		foundUser, err := userDao.FindUserByID(ctx, int(registered.ID))

		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, registered.UserName, foundUser.UserName)
	})

	t.Run("find non-existent user by id", func(t *testing.T) {
		foundUser, err := userDao.FindUserByID(ctx, 99999)

		assert.Error(t, err)
		assert.NotNil(t, foundUser)
	})
}

func TestUserDao_FindUserByUUID(t *testing.T) {
	ctx, db, cleanup := setupTestContainer(t)
	defer cleanup()

	userDao := NewUserDao(db)

	user := system.SysUser{
		UserName: "findbyuuid_user",
		Password: "password123",
		NickName: "Find By UUID Test",
	}
	registered, err := userDao.Register(ctx, user)
	assert.NoError(t, err)

	t.Run("find user by uuid string", func(t *testing.T) {
		foundUser, err := userDao.FindUserByUUID(ctx, registered.UUID.String())

		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, registered.UserName, foundUser.UserName)
	})

	t.Run("find non-existent user by uuid string", func(t *testing.T) {
		_, err := userDao.FindUserByUUID(ctx, "00000000-0000-0000-0000-000000000000")

		assert.Error(t, err)
	})
}

func TestUserDao_ResetPassword(t *testing.T) {
	ctx, db, cleanup := setupTestContainer(t)
	defer cleanup()

	userDao := NewUserDao(db)

	user := system.SysUser{
		UserName: "resetpwd_user",
		Password: "password123",
		NickName: "Reset Password Test",
	}
	registered, err := userDao.Register(ctx, user)
	assert.NoError(t, err)

	t.Run("reset password successfully", func(t *testing.T) {
		err := userDao.ResetPassword(ctx, registered.ID)

		assert.NoError(t, err)

		var updated system.SysUser
		err = db.First(&updated, registered.ID).Error
		assert.NoError(t, err)
		assert.NotEqual(t, registered.Password, updated.Password)
	})
}

func TestUserDao_UpdateLoginTime(t *testing.T) {
	ctx, db, cleanup := setupTestContainer(t)
	defer cleanup()

	userDao := NewUserDao(db)

	user := system.SysUser{
		UserName: "login_time_user",
		Password: "password123",
		NickName: "Update Login Time Test",
	}
	registered, err := userDao.Register(ctx, user)
	assert.NoError(t, err)

	t.Run("update login time successfully", func(t *testing.T) {
		err := userDao.UpdateLoginTime(ctx, registered.ID)

		assert.NoError(t, err)

		var updated system.SysUser
		err = db.First(&updated, registered.ID).Error
		assert.NoError(t, err)
		assert.NotNil(t, updated.LastLogin)
	})
}
