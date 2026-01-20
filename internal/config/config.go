package config

import "gin-admin/pkg/logging"

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
}
