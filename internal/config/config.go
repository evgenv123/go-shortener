package config

import (
	"errors"
	"net/url"
	"os"

	"github.com/alexflint/go-arg"
)

type Config struct {
	ServerAddr  string `arg:"-a,--,env:SERVER_ADDRESS" help:"Server address" default:"localhost:8080"`
	BaseURL     string `arg:"-b,--,env:BASE_URL" help:"Base URL" default:"http://localhost:8080"`
	FileStorage string `arg:"-f,--,env:FILE_STORAGE_PATH" help:"Storage filename" default:"urlStorage.gob"`
	DBSource    string `arg:"-d,--,env:DATABASE_DSN" help:"Data Source Name" default:"postgres://postgres:tttest@localhost:5432/postgres"`
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
	// Checking flags & ENV
	arg.MustParse(&conf)

	return conf, conf.Validate()
}
