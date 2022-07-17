package database

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"log"
	"sync"

	"goqrs/envs"

	"github.com/ksaucedo002/answer/errores"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type key int

var (
	dbKey key
	conn  *gorm.DB
	once  sync.Once
)

func logMode() logger.LogLevel {
	value := envs.FindEnv("GOQRS_DB_LOGS", "silent")
	switch value {
	case "info":
		return logger.Info
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	}
	return logger.Silent
}
func PrepareConnection() (err error) {
	once.Do(func() {
		host := envs.FindEnv("GOQRS_DB_HOST", "localhost")
		port := envs.FindEnv("GOQRS_DB_PORT", "5432")
		user := envs.FindEnv("GOQRS_DB_USER", "postgres")
		password := envs.FindEnv("GOQRS_DB_PASSWORD", "root")
		dbname := envs.FindEnv("GOQRS_DB_NAME", "tickets_system_db")

		const layer = "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable"
		dsn := fmt.Sprintf(layer, host, user, password, dbname, port)

		conn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			SkipDefaultTransaction: true,
			Logger:                 logger.Default.LogMode(logMode()),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
	})
	return err
}

func GormMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := WithConnection(c.Request().Context())
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
func WithConnection(ctx context.Context) context.Context {
	return context.WithValue(ctx, dbKey, conn.WithContext(ctx))
}
func Conn(ctx context.Context) *gorm.DB {
	value := ctx.Value(dbKey)
	if value == nil {
		panic("connection value not found with dbKey")
	}
	connection, ok := value.(*gorm.DB)
	if !ok {
		panic("connection invalid type")
	}
	return connection
}
func Transaction(ctx context.Context, f func(tx *gorm.DB) error) (err error) {
	tx := Conn(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			if err := tx.Rollback(); err != nil {
				log.Println("rollback error")
			}
			err = errors.New(fmt.Sprint(r))
		}
	}()
	if err := f(tx); err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return errores.NewInternalDBf(err)
	}
	return nil
}
