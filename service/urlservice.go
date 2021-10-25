package service

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/evgenv123/go-shortener/model"
	"github.com/evgenv123/go-shortener/storage"
	"runtime"
	"time"
)

type (
	Processor struct {
		config     Config
		urlStorage StorageExpected
		// Context for workers
		workerCtx context.Context
		// Cancel function for workers
		workerCancelF context.CancelFunc
		// Channel to send items to delete
		deleteCh chan model.ShortenedURL
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
	if err := c.Validate(); err != nil {
		return nil, err
	}
	if c.HexSecret == "" {
		c.HexSecret = hex.EncodeToString(defaultSecret)
	}
	if c.ChanBuffer == 0 {
		c.ChanBuffer = defaultChanBuffer
	}
	if c.WorkerThreads == 0 {
		c.WorkerThreads = runtime.NumCPU()
	}
	proc := Processor{config: c, urlStorage: st}
	// Defining channel for delete workers to receive input
	proc.deleteCh = make(chan model.ShortenedURL, c.ChanBuffer)
	// Defining context with cancel func for asynchronous worker
	proc.workerCtx, proc.workerCancelF = context.WithCancel(context.Background())
	// Starting workers
	for i := 0; i < c.WorkerThreads; i++ {
		// Making workers out of sync
		go proc.DeleteWorker(time.NewTicker(c.WorkerFlushTimeout+time.Second*time.Duration(i)), i)
	}
	return &proc, nil
}

func (svc *Processor) Close() error {
	// Stopping workers
	svc.workerCancelF()
	if svc.urlStorage == nil {
		return nil
	}
	if err := svc.urlStorage.Close(); err != nil {
		return err
	}
	return nil
}
