package app

import (
	"compress/gzip"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
)

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
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
		r.Body = gz

		next.ServeHTTP(w, r)
	})
}

// CheckSessionCookies checks if user is authorized and assigns id if not
func CheckSessionCookies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userid, err1 := r.Cookie("userid")
		sha, err2 := r.Cookie("userid-sha")
		// If we don't have cookie, or we have wrong cookie we have to set it
		if err1 != nil || err2 != nil || !checkValidAuth(userid.Value, sha.Value) {
			id, _ := uuid.NewRandom()
			cookie1 := &http.Cookie{
				Name:  "userid",
				Value: id.String(),
			}

			cookie2 := &http.Cookie{
				Name:  "userid-sha",
				Value: generateSha(id.String()),
			}
			http.SetCookie(w, cookie1)
			http.SetCookie(w, cookie2)
		}
		next.ServeHTTP(w, r)
	})
}
