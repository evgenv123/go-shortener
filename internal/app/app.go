package app

import (
	"github.com/go-chi/chi/v5"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

var m = make(map[int]string)

// MyHandlerGetId is for getting full URL from shortened
func MyHandlerGetID(w http.ResponseWriter, r *http.Request) {
	requestedID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || m[requestedID] == "" {
		http.Error(w, "Wrong requested ID!", http.StatusBadRequest)
	} else {
		http.Redirect(w, r, m[requestedID], http.StatusTemporaryRedirect)
	}
}

// MyHandlerPost is for shortening full URL and saving info to DB
func MyHandlerPost(w http.ResponseWriter, r *http.Request) {
	// reading original link body
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read request body!", http.StatusInternalServerError)
		return
	}
	// Checking for valid URL
	_, err = url.ParseRequestURI(string(b))
	if err != nil {
		http.Error(w, "Wrong URL format!", http.StatusBadRequest)
		return
	}
	// Generating ID for link (b)
	idForLink := rand.Intn(999999)
	m[idForLink] = string(b)

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("http://localhost:8080/" + strconv.Itoa(idForLink)))
	if err != nil {
		http.Error(w, "Cannot write reply body!", http.StatusInternalServerError)
		return
	}
}
