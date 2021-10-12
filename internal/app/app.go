package app

import (
	"github.com/evgenv123/go-shortener/internal/config"
)

func Init(c config.Config) error {
	var err error
	appConf = c
	URLSvc, err = c.BuildURLService()
	return err
}

func Close() error {
	return URLSvc.Close()
}
