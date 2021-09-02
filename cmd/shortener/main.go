package main

import (
	"github.com/evgenv123/go-shortener/internal/app"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	// маршрутизация запросов обработчику
	r.Get("/{id}", app.MyHandler)
	r.Post("/", app.MyHandler)
	// запуск сервера с адресом localhost, порт 8080
	log.Fatal(http.ListenAndServe(":8080", r))
}
