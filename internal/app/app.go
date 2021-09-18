package app

import (
	"github.com/evgenv123/go-shortener/internal/config"
)

func Init(c config.Config) error {
	appConf = c
	if err := readDBFromFile(); err != nil {
		return err
	}
	return nil
}
