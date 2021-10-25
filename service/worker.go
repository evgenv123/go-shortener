package service

import (
	"context"
	"github.com/evgenv123/go-shortener/model"
	"log"
	"time"
)

// DeleteWorker is asynchronous goroutine
func (svc *Processor) DeleteWorker(ticker *time.Ticker, workerID int) {
	var buffer []model.ShortenedURL
	for {
		select {
		case item := <-svc.deleteCh:
			buffer = append(buffer, item)
		case <-ticker.C:
			log.Printf("Worker #%v:  flushing buffer\n", workerID)
			// Flushing after worker timeout
			if err := svc.urlStorage.DeleteBatchURL(context.Background(), buffer); err != nil {
				log.Printf("Worker #%v:  error batch deleting urls\n", workerID)
			}
			buffer = nil
		case <-svc.workerCtx.Done():
			if err := svc.urlStorage.DeleteBatchURL(context.Background(), buffer); err != nil {
				log.Printf("Worker #%v:  error batch deleting urls\n", workerID)
			}
			buffer = nil
			log.Printf("Worker #%v:  gracefully stopped DeleteWorker\n", workerID)
			return
		}
	}
}
