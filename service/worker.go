package service

import (
	"context"
	"github.com/evgenv123/go-shortener/model"
	"log"
	"time"
)

// DeleteWorker is asynchronous goroutine
func (svc *Processor) DeleteWorker(tickerCh <-chan time.Time, workerId int) {
	var buffer []model.ShortenedURL
	for {
		select {
		case item := <-svc.deleteCh:
			buffer = append(buffer, item)
		case <-tickerCh:
			log.Printf("Worker #%v:  flushing buffer\n", workerId)
			// Flushing after worker timeout
			if err := svc.urlStorage.DeleteBatchURL(context.Background(), buffer); err != nil {
				log.Printf("Worker #%v:  error batch deleting urls\n", workerId)
			}
			buffer = nil
		case <-svc.workerCtx.Done():
			if err := svc.urlStorage.DeleteBatchURL(context.Background(), buffer); err != nil {
				log.Printf("Worker #%v:  error batch deleting urls\n", workerId)
			}
			buffer = nil
			log.Printf("Worker #%v:  gracefully stopped DeleteWorker\n", workerId)
			return
		}
	}
}
