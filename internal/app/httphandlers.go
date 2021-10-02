package app

import (
	"encoding/json"
	"errors"
	"github.com/evgenv123/go-shortener/internal/dbcore"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// MyHandlerGetID is for getting full URL from shortened
func MyHandlerGetID(w http.ResponseWriter, r *http.Request) {
	// TODO: Read from SQL if online
	requestedID, err := strconv.Atoi(chi.URLParam(r, "id"))
	DB.RLock()
	if err != nil || DB.URLMap[requestedID].URL == "" {
		http.Error(w, "Wrong requested ID!", http.StatusBadRequest)
	} else {
		http.Redirect(w, r, DB.URLMap[requestedID].URL, http.StatusTemporaryRedirect)
	}
	DB.RUnlock()
}

// MyHandlerListUrls is for getting all URLS for specified user
func MyHandlerListUrls(w http.ResponseWriter, r *http.Request) {
	// TODO: Read from SQL if online
	var result []OutputAllURLs
	DB.RLock()
	// Iterating over all URLs
	for k, v := range DB.URLMap {
		// Appending to result if it matches our username
		if v.UserID == r.Context().Value(contextKeyUserID).(string) {
			result = append(result, OutputAllURLs{ShortURL: GetShortenedURL(k), OriginalURL: v.URL})
		}
	}
	DB.RUnlock()

	if len(result) == 0 {
		http.Error(w, "No links for user!", http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, "Cannot write reply body!", http.StatusInternalServerError)
		return
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

	// Trying to shorten URL from request
	shortened, err := shortenURL(string(b), r.Context().Value(contextKeyUserID).(string))

	if err != nil {
		// If we receive duplicate error from SQL
		var myErr *FullURLDuplicateError
		if errors.As(err, &myErr) {
			// Send StatusConflict and existing short url
			w.WriteHeader(http.StatusConflict)
			shortened = myErr.ShortURL
		} else {
			http.Error(w, "Cannot shorten URL!"+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	// Writing reply
	_, err = w.Write([]byte(shortened))
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

	shortened, err := shortenURL(input.URL, r.Context().Value(contextKeyUserID).(string))
	if err != nil {
		// If we receive duplicate error from SQL
		var myErr *FullURLDuplicateError
		if errors.As(err, &myErr) {
			// Send StatusConflict and existing short url
			w.WriteHeader(http.StatusConflict)
			shortened = myErr.ShortURL
		} else {
			http.Error(w, "Cannot shorten URL!"+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	err = json.NewEncoder(w).Encode(OutputShortURL{Result: shortened})
	if err != nil {
		http.Error(w, "Cannot write reply body!", http.StatusInternalServerError)
		return
	}
}

// MyHandlerShortenBatch is a handler for /api/shorten/batch endpoint
func MyHandlerShortenBatch(w http.ResponseWriter, r *http.Request) {
	// reading original link body
	var input []InputBatch
	var output []OutputBatch
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Cannot decode request body!", http.StatusInternalServerError)
		return
	}
	for i := 0; i < len(input); i++ {
		// Checking for valid URL
		_, err := url.ParseRequestURI(input[i].OrigURL)
		if err != nil {
			http.Error(w, "Wrong URL format!", http.StatusBadRequest)
			return
		}
		shortened, _ := shortenURL(input[i].OrigURL, r.Context().Value(contextKeyUserID).(string))
		output = append(output, OutputBatch{
			CorrID:   input[i].CorrID,
			ShortURL: shortened,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, "Cannot write reply body!", http.StatusInternalServerError)
		return
	}
}

// MyHandlerPing is a handler for /ping endpoint
func MyHandlerPing(w http.ResponseWriter, r *http.Request) {
	if !dbcore.CheckConn() {
		http.Error(w, "Cannot ping database!", http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
