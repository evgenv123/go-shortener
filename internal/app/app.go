package app

import (
	"encoding/json"
	"github.com/evgenv123/go-shortener/internal/config"
	"github.com/go-chi/chi/v5"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

var m = make(map[int]string)

type InputURL struct {
	URL string `json:"url"`
}
type OutputShortURL struct {
	Result string `json:"result"`
}

func shortenURL(url string) OutputShortURL {
	// Generating ID for link (b)
	idForLink := rand.Intn(999999)
	m[idForLink] = url
	return OutputShortURL{Result: config.BaseURL + "/" + strconv.Itoa(idForLink)}
}

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
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortenURL(string(b)).Result))
	if err != nil {
		http.Error(w, "Cannot write reply body!", http.StatusInternalServerError)
		return
	}
}

// MyHandlerShorten is a handler for /api/shorten endpoint
func MyHandlerShorten(w http.ResponseWriter, r *http.Request) {
	// reading original link body
	var input InputURL
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Cannot decode request body!", http.StatusInternalServerError)
		return
	}
	// Checking for valid URL
	_, err := url.ParseRequestURI(input.URL)
	if err != nil {
		http.Error(w, "Wrong URL format!", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(shortenURL(input.URL))
	if err != nil {
		http.Error(w, "Cannot write reply body!", http.StatusInternalServerError)
		return
	}
}
