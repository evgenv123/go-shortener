package app

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// MyHandlerGetId is for getting full URL from shortened
func MyHandlerGetID(w http.ResponseWriter, r *http.Request) {
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
	var result []OutputAllURLs
	// no error checking on cookie because we have middleware to check cookies
	userid, _ := r.Cookie("userid")
	DB.RLock()
	// Iterating over all URLs
	for k, v := range DB.URLMap {
		// Appending to result if it matches our username
		if v.UserID == userid.Value {
			result = append(result, OutputAllURLs{ShortURL: getShortenedURL(k), OriginalURL: v.URL})
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
	w.WriteHeader(http.StatusCreated)
	// no error checking on cookie because we have middleware to check cookies
	userid, _ := r.Cookie("userid")
	_, err = w.Write([]byte(shortenURL(string(b), userid.Value).Result))
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
	// no error checking on cookie because we have middleware to check cookies
	userid, _ := r.Cookie("userid")
	err = json.NewEncoder(w).Encode(shortenURL(input.URL, userid.Value))
	if err != nil {
		http.Error(w, "Cannot write reply body!", http.StatusInternalServerError)
		return
	}
}
