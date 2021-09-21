package app

import (
	"compress/gzip"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w gzipReader) Read(b []byte) (int, error) {
	return w.Reader.Read(b)
}

func GZipWriteHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func GZipReadHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент отправил сжатый gzip-запрос
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Reader
		gz, err := gzip.NewReader(r.Body)

		if err != nil && err != io.EOF {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		next.ServeHTTP(w, r)
	})
}

// MyHandlerGetId is for getting full URL from shortened
func MyHandlerGetID(w http.ResponseWriter, r *http.Request) {
	requestedID, err := strconv.Atoi(chi.URLParam(r, "id"))
	DB.RLock()
	if err != nil || DB.URLMap[requestedID] == "" {
		http.Error(w, "Wrong requested ID!", http.StatusBadRequest)
	} else {
		http.Redirect(w, r, DB.URLMap[requestedID], http.StatusTemporaryRedirect)
	}
	DB.RUnlock()
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
