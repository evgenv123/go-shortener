package config

import (
	"errors"
	"net/url"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddr  string `arg:"-a" help:"Server address" env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL     string `arg:"-b" help:"Base URL" env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStorage string `arg:"-f" help:"Storage filename" env:"FILE_STORAGE_PATH" envDefault:"urlStorage.gob"`
}

func (c Config) Validate() error {
	_, err := os.Stat(c.FileStorage)
	if err != nil && !os.IsNotExist(err) {
		return errors.New("wrong file name")
	}
	_, err = url.Parse(c.ServerAddr)
	if err != nil {
		return errors.New("wrong server address")
	}
	_, err = url.ParseRequestURI(c.BaseURL)
	if err != nil {
		return errors.New("wrong base url")
	}
	return nil
}

func NewConfig() (Config, error) {
	var conf Config
	// Checking flags first, then redefining them by ENV
	arg.MustParse(&conf)
	err := env.Parse(&conf)
	if err != nil {
		return Config{
			"localhost:8080",
			"http://localhost:8080",
			"urlStorage.gob",
		}, nil
	}
	return conf, conf.Validate()
}
