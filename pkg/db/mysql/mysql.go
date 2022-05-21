package mysql

import (
	"context"
	sysLog "log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/yusank/goim/pkg/graceful"
	"github.com/yusank/goim/pkg/log"
)

var (
	// default db
	defaultDB *gorm.DB
)

func GetDB() *gorm.DB {
	return defaultDB
}

func InitDB(opts ...Option) error {
	gdb, err := NewMySQL(opts...)
	if err != nil {
		return err
	}

	SetDefaultDB(gdb)
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
		SlowThreshold: 1 * time.Second,
		LogLevel:      logger.Warn,
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
		Logger: slowLogger,
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

// SetDefaultDB set default db, If you want to use default db, you can just use GetDB()
// If called InitDB, you don't need to call this function.
func SetDefaultDB(db *gorm.DB) {
	defaultDB = db
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
