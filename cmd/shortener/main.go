package main

import (
	"github.com/evgenv123/go-shortener/internal/app"
	"github.com/evgenv123/go-shortener/internal/config"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func init() {
	config.Init()
}

func main() {
	r := chi.NewRouter()
	// маршрутизация запросов обработчику
	r.Get("/{id}", app.MyHandlerGetID)
	r.Post("/api/shorten", app.MyHandlerShorten)
	r.Post("/", app.MyHandlerPost)
	log.Fatal(http.ListenAndServe(config.ServerAddr, r))
}
