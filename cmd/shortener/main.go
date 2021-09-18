package main

import (
	"context"
	"github.com/evgenv123/go-shortener/internal/app"
	"github.com/evgenv123/go-shortener/internal/config"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func startHTTP(srv *http.Server) {
	// Error ErrServerClosed is thrown during graceful shutdown
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	if err := app.Init(conf); err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	// маршрутизация запросов обработчику
	r.Get("/{id}", app.MyHandlerGetID)
	r.Post("/api/shorten", app.MyHandlerShorten)
	r.Post("/", app.MyHandlerPost)
	srv := &http.Server{
		Addr:    conf.ServerAddr,
		Handler: r,
	}
	// Creating interrupt channel
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	// Starting web server in background
	go startHTTP(srv)
	// Waiting signal for shutdown
	<-done
	// Giving server 5 sec to shut down
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown error: %v\n", err)
	} else {
		log.Printf("gracefully stopped\n")
	}
	// Writing DB to file on exit
	if err = app.WriteDBToFile(); err != nil {
		log.Println("Error writing DB to file: ", err)
	}
}
