package main

import (
	"os"

	"github.com/BurntSushi/toml"
	grpcservice "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/server/grpcservice"
	internalhttp "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/server/http"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	HTTPServer internalhttp.Conf `toml:"http"`
	GRPSServer grpcservice.Conf  `toml:"grpc"`
	Storage    StorageConf       `toml:"storage"`
	Logger     LoggerConf        `toml:"logger"`
}

type LoggerConf struct {
	Level string `toml:"level"`
}

type StorageConf struct {
	DB  string `toml:"db"`
	DSN string `toml:"dsn"`
}

func NewConfig() Config {
	return Config{}
}

func (c *Config) LoadFileTOML(filename string) error {
	filedata, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return toml.Unmarshal(filedata, c)
}
