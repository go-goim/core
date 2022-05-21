package mysql

import (
	"time"

	configv1 "github.com/yusank/goim/api/config/v1"
)

type Option func(*option)

type option struct {
	addr            string
	user            string
	password        string
	db              string
	dsn             string
	maxConns        int
	maxIdleConns    int
	idleTimeout     time.Duration
	connMaxLifetime time.Duration
	debug           bool
}

func newOption() *option {
	return &option{
		addr:            "127.0.0.1:3306",
		user:            "goim",
		db:              "goim",
		dsn:             "",
		maxConns:        100,
		maxIdleConns:    10,
		idleTimeout:     time.Minute,
		connMaxLifetime: time.Minute * 10,
	}
}

func (o *option) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func Addr(addr string) Option {
	return func(o *option) {
		o.addr = addr
	}
}

func User(user string) Option {
	return func(o *option) {
		o.user = user
	}
}

func Password(password string) Option {
	return func(o *option) {
		o.password = password
	}
}

func DBName(db string) Option {
	return func(o *option) {
		o.db = db
	}
}

func DSN(dsn string) Option {
	return func(o *option) {
		o.dsn = dsn
	}
}

func MaxConns(i int) Option {
	return func(o *option) {
		o.maxConns = i
	}
}

func MaxIdleConns(i int) Option {
	return func(o *option) {
		o.maxIdleConns = i
	}
}

func IdleTimeout(d time.Duration) Option {
	return func(o *option) {
		o.idleTimeout = d
	}
}

func ConnMaxLifetime(d time.Duration) Option {
	return func(o *option) {
		o.connMaxLifetime = d
	}
}

func Debug(b bool) Option {
	return func(o *option) {
		o.debug = b
	}
}

func WithConfig(cfg *configv1.MySQL) Option {
	return func(o *option) {
		o.addr = cfg.GetAddr()
		o.user = cfg.GetUser()
		o.password = cfg.GetPassword()
		o.db = cfg.GetDb()
		o.maxConns = int(cfg.GetMaxOpenConns())
		o.maxIdleConns = int(cfg.GetMaxIdleConns())

		idleTimeout := cfg.GetIdleTimeout()
		if idleTimeout != nil {
			o.idleTimeout = idleTimeout.AsDuration()
		}

		connMaxLifetime := cfg.GetOpenTimeout()
		if connMaxLifetime != nil {
			o.connMaxLifetime = connMaxLifetime.AsDuration()
		}
	}
}
