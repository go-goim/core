package mysql

import (
	"context"
	sysLog "log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/go-goim/core/pkg/graceful"
	"github.com/go-goim/core/pkg/log"
)

var (
	// default db
	defaultDB *gorm.DB
)

func GetDB() *gorm.DB {
	return defaultDB
}

func InitDB(opts ...Option) error {
	var err error
	defaultDB, err = NewMySQL(opts...)
	if err != nil {
		return err
	}

	graceful.Register(func(ctx context.Context) error {
		err := Close()
		if err != nil {
			log.Error("mysql close error", "err", err)
		}

		return err
	})
	return nil
}

func NewMySQL(opts ...Option) (*gorm.DB, error) {
	o := newOption()
	o.apply(opts...)

	if o.dsn == "" {
		o.dsn = o.user + ":" + o.password + "@tcp(" + o.addr + ")/" + o.db + "?charset=utf8mb4&parseTime=True&loc=Local"
	}

	loggerConfig := logger.Config{
		SlowThreshold:             1 * time.Second,
		LogLevel:                  logger.Error,
		IgnoreRecordNotFoundError: true,
	}

	if o.debug {
		loggerConfig = logger.Config{
			LogLevel: logger.Info,
		}
	}

	slowLogger := logger.New(
		sysLog.New(os.Stdout, "[MYSQL]", sysLog.LstdFlags),
		loggerConfig,
	)

	gConf := &gorm.Config{
		Logger:                 slowLogger,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	gdb, err := gorm.Open(mysql.Open(o.dsn), gConf)
	if err != nil {
		return nil, err
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(o.maxIdleConns)
	sqlDB.SetMaxOpenConns(o.maxConns)
	sqlDB.SetConnMaxIdleTime(o.idleTimeout)
	sqlDB.SetConnMaxLifetime(o.connMaxLifetime)

	return gdb, nil
}

// Close db.
func Close() error {
	if defaultDB != nil {
		db, err := defaultDB.DB()
		if err != nil {
			return err
		}
		return db.Close()
	}

	return nil
}
