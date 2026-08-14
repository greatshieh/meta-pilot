package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"server/pkg/model/common/request"
	"server/pkg/model/system"
)

func setupAPITestContainer(t *testing.T) (context.Context, *gorm.DB, func()) {
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

	err = db.AutoMigrate(&system.SysApi{})
	assert.NoError(t, err)

	cleanup := func() {
		mysqlC.Terminate(ctx)
	}

	return ctx, db, cleanup
}

func TestApiDao_CreateAPI(t *testing.T) {
	ctx, db, cleanup := setupAPITestContainer(t)
	defer cleanup()

	apiDao := NewApiDao(db)

	t.Run("create new api", func(t *testing.T) {
		api := system.SysApi{
			Path:        "/api/test",
			Method:      "GET",
			Description: "Test API",
			ApiGroup:    "test",
		}

		err := apiDao.CreateAPI(ctx, api)
		assert.NoError(t, err)
	})

	t.Run("create existing api", func(t *testing.T) {
		api := system.SysApi{
			Path:        "/api/test",
			Method:      "GET",
			Description: "Test API",
			ApiGroup:    "test",
		}

		err := apiDao.CreateAPI(ctx, api)
		assert.Error(t, err)
	})
}

func TestApiDao_DeleteAPI(t *testing.T) {
	ctx, db, cleanup := setupAPITestContainer(t)
	defer cleanup()

	apiDao := NewApiDao(db)

	api := system.SysApi{
		Path:        "/api/delete",
		Method:      "DELETE",
		Description: "Delete API",
		ApiGroup:    "test",
	}
	err := db.Create(&api).Error
	assert.NoError(t, err)

	t.Run("delete existing api", func(t *testing.T) {
		deletedApi, err := apiDao.DeleteAPI(ctx, api.ID)
		assert.NoError(t, err)
		assert.Equal(t, api.ID, deletedApi.ID)
	})

	t.Run("delete non-existent api", func(t *testing.T) {
		_, err := apiDao.DeleteAPI(ctx, 99999)
		assert.Error(t, err)
	})
}

func TestApiDao_GetAPIInfoList(t *testing.T) {
	ctx, db, cleanup := setupAPITestContainer(t)
	defer cleanup()

	apiDao := NewApiDao(db)

	for i := 0; i < 5; i++ {
		api := system.SysApi{
			Path:        "/api/test" + string(rune('a'+i)),
			Method:      "GET",
			Description: "Test API " + string(rune('a'+i)),
			ApiGroup:    "test",
		}
		err := db.Create(&api).Error
		assert.NoError(t, err)
	}

	t.Run("get api list with pagination", func(t *testing.T) {
		pageInfo := request.PageInfo{
			Page:     1,
			PageSize: 2,
		}

		list, total, err := apiDao.GetAPIInfoList(ctx, system.SysApi{}, pageInfo)
		assert.NoError(t, err)
		assert.Len(t, list, 2)
		assert.Equal(t, int64(5), total)
	})

	t.Run("get api list with path filter", func(t *testing.T) {
		pageInfo := request.PageInfo{
			Page:     1,
			PageSize: 10,
		}

		list, _, err := apiDao.GetAPIInfoList(ctx, system.SysApi{Path: "testa"}, pageInfo)
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})

	t.Run("get api list with group filter", func(t *testing.T) {
		pageInfo := request.PageInfo{
			Page:     1,
			PageSize: 10,
		}

		list, _, err := apiDao.GetAPIInfoList(ctx, system.SysApi{ApiGroup: "test"}, pageInfo)
		assert.NoError(t, err)
		assert.Len(t, list, 5)
	})
}

func TestApiDao_UpdateAPI(t *testing.T) {
	ctx, db, cleanup := setupAPITestContainer(t)
	defer cleanup()

	apiDao := NewApiDao(db)

	api := system.SysApi{
		Path:        "/api/update",
		Method:      "PUT",
		Description: "Update API",
		ApiGroup:    "test",
	}
	err := db.Create(&api).Error
	assert.NoError(t, err)

	t.Run("update api successfully", func(t *testing.T) {
		api.Description = "Updated API"
		err := apiDao.UpdateAPI(ctx, api)
		assert.NoError(t, err)

		var updatedApi system.SysApi
		err = db.First(&updatedApi, api.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, "Updated API", updatedApi.Description)
	})

	t.Run("update api with duplicate path and method", func(t *testing.T) {
		anotherApi := system.SysApi{
			Path:        "/api/another",
			Method:      "PUT",
			Description: "Another API",
			ApiGroup:    "test",
		}
		err := db.Create(&anotherApi).Error
		assert.NoError(t, err)

		anotherApi.Path = "/api/update"
		err = apiDao.UpdateAPI(ctx, anotherApi)
		assert.Error(t, err)
	})
}

func TestApiDao_GetAllAPIs(t *testing.T) {
	ctx, db, cleanup := setupAPITestContainer(t)
	defer cleanup()

	apiDao := NewApiDao(db)

	for i := 0; i < 3; i++ {
		api := system.SysApi{
			Path:        "/api/all" + string(rune('a'+i)),
			Method:      "GET",
			Description: "All API " + string(rune('a'+i)),
			ApiGroup:    "all",
		}
		err := db.Create(&api).Error
		assert.NoError(t, err)
	}

	t.Run("get all apis", func(t *testing.T) {
		apis, err := apiDao.GetAllAPIs(ctx, 1)
		assert.NoError(t, err)
		assert.Len(t, apis, 3)
	})
}

func TestApiDao_GetAPIGroups(t *testing.T) {
	ctx, db, cleanup := setupAPITestContainer(t)
	defer cleanup()

	apiDao := NewApiDao(db)

	for i := 0; i < 3; i++ {
		api := system.SysApi{
			Path:        "/api/group" + string(rune('a'+i)),
			Method:      "GET",
			Description: "Group API " + string(rune('a'+i)),
			ApiGroup:    "group" + string(rune('a'+i)),
		}
		err := db.Create(&api).Error
		assert.NoError(t, err)
	}

	t.Run("get api groups", func(t *testing.T) {
		groups, err := apiDao.GetAPIGroups(ctx)
		assert.NoError(t, err)
		assert.Len(t, groups, 3)
	})
}
