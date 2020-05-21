package internal

import (
	"github.com/anastasja-hunko/test/internal/database"
)

type Config struct {
	Port     string
	LogLevel string
	DbConfig *database.Config
}

func NewConfig() *Config {
	return &Config{
		Port:     ":8181",
		LogLevel: "debug",
		DbConfig: database.NewConfig(),
	}
}
