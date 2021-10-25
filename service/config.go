package service

import (
	"encoding/hex"
	"errors"
	"net/url"
	"time"
)

type Config struct {
	BaseURL            string
	HexSecret          string
	WorkerFlushTimeout time.Duration
	WorkerThreads      int
	ChanBuffer         int
}

var (
	defaultSecret     = []byte{18, 232, 139, 12, 216, 189, 22, 128, 122, 49, 246, 137, 191, 24, 38, 210}
	defaultChanBuffer = 20
)

func (c Config) Validate() error {
	_, err := url.ParseRequestURI(c.BaseURL)
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
	if c.WorkerFlushTimeout < time.Second*1 || c.WorkerFlushTimeout > time.Second*100 {
		return errors.New("worker flush timeout should be in range 1s..100s")
	}
	if c.ChanBuffer != 0 && (c.ChanBuffer < 1 || c.ChanBuffer > 100) {
		return errors.New("ChanBuffer should be in range 1..100")
	}
	if c.WorkerThreads != 0 && (c.WorkerThreads < 1 || c.WorkerThreads > 100) {
		return errors.New("WorkerThreads should be in range 1..100")
	}
	return nil
}
