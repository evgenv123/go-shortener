package main

import (
	"context"
	"encoding/json"
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
	// Error ErrServerClosed is thrown during graceful shutdown, so we consider it is not error
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Listen: %s\n", err)
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
	r.Get("/user/urls", app.MyHandlerListUrls)
	r.Post("/api/shorten", app.MyHandlerShorten)
	r.Post("/", app.MyHandlerPost)
	srv := &http.Server{
		Addr:    conf.ServerAddr,
		Handler: app.Conveyor(r, app.CheckSessionCookies, app.GZipReadHandler, app.GZipWriteHandler),
	}
	// Creating interrupt channel
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	// Starting web server in background
	jsonConf, _ := json.Marshal(conf)
	log.Println("Starting web server:", string(jsonConf))
	go startHTTP(srv)
	// Waiting signal for shutdown
	<-done
	// Giving server 5 sec to shut down
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v\n", err)
	} else {
		log.Printf("Gracefully stopped\n")
	}
	// Writing DB to file on exit
	if err = app.WriteDBToFile(); err != nil {
		log.Println("Error writing DB to file: ", err)
	}
}
