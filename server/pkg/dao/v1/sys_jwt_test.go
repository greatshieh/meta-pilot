package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"server/pkg/model/system"
)

func setupJwtTestContainer(t *testing.T) (context.Context, *gorm.DB, func()) {
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

	err = db.AutoMigrate(&system.JwtBlacklist{})
	assert.NoError(t, err)

	cleanup := func() {
		mysqlC.Terminate(ctx)
	}

	return ctx, db, cleanup
}

func TestJwtDao_JsonInBlacklist(t *testing.T) {
	ctx, db, cleanup := setupJwtTestContainer(t)
	defer cleanup()

	jwtDao := NewJwtDao(db)

	t.Run("add jwt to blacklist", func(t *testing.T) {
		jwt := system.JwtBlacklist{
			Jwt: "test-jwt-token",
		}

		err := jwtDao.JsonInBlacklist(ctx, jwt)
		assert.NoError(t, err)
	})
}

func TestJwtDao_IsBlacklist(t *testing.T) {
	ctx, db, cleanup := setupJwtTestContainer(t)
	defer cleanup()

	jwtDao := NewJwtDao(db)

	t.Run("check non-existent jwt", func(t *testing.T) {
		isBlacklisted := jwtDao.IsBlacklist(ctx, "non-existent-token")
		assert.False(t, isBlacklisted)
	})

	t.Run("check existing jwt in cache", func(t *testing.T) {
		jwt := system.JwtBlacklist{
			Jwt: "test-jwt-token-cache",
		}
		err := jwtDao.JsonInBlacklist(ctx, jwt)
		assert.NoError(t, err)

		isBlacklisted := jwtDao.IsBlacklist(ctx, "test-jwt-token-cache")
		assert.True(t, isBlacklisted)
	})
}
