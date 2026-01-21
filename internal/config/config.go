package config

import (
	"encoding/json"
	"fmt"
	"gin-admin/pkg/logging"
)

type Config struct {
	Logger     logging.LoggerConfig
	General    General
	Storage    Storage
	Middleware Middleware
	Util       Util
	Dictionary Dictionary
}

type General struct {
	AppName            string
	Version            string
	Debug              bool
	PprofAddr          string
	DisableSwagger     bool
	DisablePrintConfig bool
	DefaultLoginPwd    string
	WorkDir            string
	MenuFile           string
	DenyOperateMenu    bool
	HTTP               struct {
		Addr            string
		ShutdownTimeout int
		ReadTimeout     int
		WriteTimeout    int
		IdleTimeout     int
		CertFile        string
		KeyFile         string
	}
	Root struct {
		ID       string
		Username string
		Password string
		Name     string
	}
}

type Storage struct {
	Cache struct {
		Type      string
		Delimiter string
		Memory    struct {
			CleanupInterval int
		}
		Badger struct {
			Path string
		}
		Redis struct {
			Addr     string
			Username string
			Password string
			DB       int
		}
	}
	DB struct {
		Debug        bool
		Type         string
		DSN          string
		MaxLifetime  int
		MaxIdleTime  int
		MaxOpenConns int
		TablePrefix  string
		AutoMigrate  bool
		PrepareStmt  bool
		Resolver     []struct {
			DBType   string
			Sources  []string
			Replicas []string
			Tables   []string
		}
	}
}

type Util struct {
	Captcha struct {
		Length    int
		Width     int
		Height    int
		CacheType string
		Redis     struct {
			Addr      string
			Username  string
			Password  string
			DB        int
			KeyPrefix string
		}
	}
	Prometheus struct {
		Enable         bool
		Port           int
		BasicUsername  string
		BasicPassword  string
		LogApis        []string
		LogMethods     []string
		DefaultCollect bool
	}
}

type Dictionary struct {
	UserCacheExp int
}

func (c *Config) IsDebug() bool {
	return c.General.Debug
}

func (c *Config) String() string {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic("Failed to marshal config:" + err.Error())
	}

	return string(b)
}

func (c *Config) PreLoad() {
	if addr := c.Storage.Cache.Redis.Addr; addr != "" {
		username := c.Storage.Cache.Redis.Username
		password := c.Storage.Cache.Redis.Password

		if c.Util.Captcha.CacheType == "redis" &&
			c.Util.Captcha.Redis.Addr == "" {
			c.Util.Captcha.Redis.Addr = addr
			c.Util.Captcha.Redis.Username = username
			c.Util.Captcha.Redis.Password = password
		}

		if c.Middleware.RateLimiter.Store.Type == "redis" &&
			c.Middleware.RateLimiter.Store, Redis.Addr == "" {
			c.Middleware.RateLimiter.Store.Redis.Addr = addr
			c.Middleware.RateLimiter.Store.Redis.Username = username
			c.Middleware.RateLimiter.Store.Redis.Password = password
		}

		if c.Middleware.Auth.Store.Type == "redis" &&
			c.Middleware.Auth.Store.Redis.Addr == "" {
			c.Middleware.Auth.Store.Redis.Addr = addr
			c.Middleware.Auth.Store.Redis.Username = username
			c.Middleware.Auth.Store.Redis.Password = password
		}
	}
}

func (c *Config) Print() {
	if c.General.DisablePrintConfig {
		return
	}
	fmt.Println("// ----------------------- Load configurations start ------------------------")
	fmt.Println(c.String())
	fmt.Println("// ----------------------- Load configurations end --------------------------")
}

func (c *Config) FormatTableName(name string) string {
	return c.Storage.DB.TablePrefix + name
}
