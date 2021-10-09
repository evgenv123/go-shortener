package app

import (
	"github.com/evgenv123/go-shortener/internal/config"
)

func Init(c config.Config) error {
	var err error
	appConf = c
	UrlSvc, err = c.BuildURLService()
	return err
}

func Close() error {
	return UrlSvc.Close()
}
