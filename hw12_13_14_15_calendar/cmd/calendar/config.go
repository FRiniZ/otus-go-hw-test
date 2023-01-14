package main

import (
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	HTTPServer HTTPServerConf `toml:"http"`
	Storage    StorageConf    `toml:"storage"`
	Logger     LoggerConf     `toml:"logger"`
}

type HTTPServerConf struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
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
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	filedata, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(filedata, c)
	if err != nil {
		return err
	}

	return nil
}
