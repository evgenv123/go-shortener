package config

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/evgenv123/go-shortener/service"
	"github.com/evgenv123/go-shortener/storage/gob"
	"github.com/evgenv123/go-shortener/storage/psql"
	"net/url"
	"os"
	"time"

	"github.com/alexflint/go-arg"
)

type Config struct {
	ServerAddr    string        `arg:"-a,--,env:SERVER_ADDRESS" help:"Server address" default:"localhost:8080"`
	BaseURL       string        `arg:"-b,--,env:BASE_URL" help:"Base URL" default:"http://localhost:8080"`
	FileStorage   string        `arg:"-f,--,env:FILE_STORAGE_PATH" help:"Storage filename" default:"urlStorage.gob"`
	DBSource      string        `arg:"-d,--,env:DATABASE_DSN" help:"Data Source Name (ex.: postgres://postgres:tttest@localhost:5432/postgres)"`
	HexSecret     string        `arg:"env:SECRET_KEY" help:"16-byte HEX encoded secret for cookie signing"`
	CtxTimeout    time.Duration `arg:"-t" help:"Context timeout for operations, seconds" default:"5s"`
	OperationMode OPMode
}

type OPMode int

const (
	OPModePGSQL = 1
	OPModeFile  = 2
)

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
	if c.HexSecret != "" {
		s, err := hex.DecodeString(c.HexSecret)
		if err != nil {
			return errors.New("cannot decode hex secret")
		}
		if len(s) != 16 {
			return errors.New("wrong secret key length")
		}
	}
	if c.CtxTimeout < time.Second*1 || c.CtxTimeout > time.Second*100 {
		return errors.New("context timeout should be in range 1s..100s")
	}
	return nil
}

func (c Config) BuildGOBStorage() (*gob.Storage, error) {
	gobConf := gob.Config{Filename: c.FileStorage}
	st, err := gob.New(gobConf)
	if err != nil {
		return nil, fmt.Errorf("error building GOB storage: %w", err)
	}
	return st, nil
}

func (c Config) BuildPSQLStorage() (psql.Storage, error) {
	psqlConf := psql.Config{DSN: c.DBSource}
	st, err := psql.New(psqlConf)
	if err != nil {
		return st, fmt.Errorf("error building PSQL storage: %w", err)
	}
	return st, nil
}

func (c Config) BuildURLService() (*service.Processor, error) {
	var st service.StorageExpected
	var err error
	// Choosing storage depending on configuration info
	if c.OperationMode == OPModeFile {
		//st = st.(gob.Storage)
		st, err = c.BuildGOBStorage()
		if err != nil {
			return nil, err
		}
	} else if c.OperationMode == OPModePGSQL {
		//st = st.(psql.Storage)
		st, err = c.BuildPSQLStorage()
		if err != nil {
			return nil, err
		}
	}
	svcConfig := service.Config{BaseURL: c.BaseURL, HexSecret: c.HexSecret}
	svc, err := service.New(svcConfig, st)
	if err != nil {
		return nil, err
	}
	return svc, nil
}

func NewConfig() (Config, error) {
	var conf Config
	// Checking flags & ENV
	arg.MustParse(&conf)

	// If any of config parameters is invalid we fail
	if err := conf.Validate(); err != nil {
		return Config{}, err
	}
	if conf.DBSource != "" {
		// SQL is our first priority
		conf.OperationMode = OPModePGSQL
	} else {
		// If no DSN specified we fallback to GOB file
		conf.OperationMode = OPModeFile
	}

	return conf, nil
}
