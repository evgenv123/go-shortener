package service

import (
	"encoding/hex"
	"errors"
	"github.com/evgenv123/go-shortener/storage"
)

type (
	Processor struct {
		config     Config
		urlStorage StorageExpected
	}
	// StorageExpected includes all types of storage urlservice can operate with
	StorageExpected interface {
		storage.URLReader
		storage.URLWriter
	}
)

func New(c Config, st StorageExpected) (*Processor, error) {
	if st == nil {
		return nil, errors.New("storage cannot be nil")
	}
	proc := Processor{config: c, urlStorage: st}
	if c.HexSecret == "" {
		proc.config.HexSecret = hex.EncodeToString(defaultSecret)
	}
	return &proc, nil
}

func (svc *Processor) Close() error {
	if svc.urlStorage == nil {
		return nil
	}
	if err := svc.urlStorage.Close(); err != nil {
		return err
	}
	return nil
}
