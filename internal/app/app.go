package app

import (
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

var m = make(map[int]string)

func MyHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && r.URL.String() == "/":
		// reading original link body
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
			return
		}
		// Generating ID for link (b)
		idForLink := rand.Intn(999999)
		m[idForLink] = string(b)

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte("http://localhost:8080/" + strconv.Itoa(idForLink)))
		if err != nil {
			log.Fatal(err)
			return
		}
	case r.Method == http.MethodGet && r.URL.String() != "/":
		requestedID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || m[requestedID] == "" {
			http.Error(w, "Wrong requested ID!", http.StatusBadRequest)
		} else {
			http.Redirect(w, r, m[requestedID], http.StatusTemporaryRedirect)
		}
	default:
		http.Error(w, "Wrong request!", http.StatusBadRequest)
	}
}
