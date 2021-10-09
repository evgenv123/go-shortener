package app

import (
	"compress/gzip"
	"context"
	"io"
	"log"
	"net/http"
	"strings"
)

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// RequestLogHandler logs all requests for debugging
func RequestLogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.String())
		next.ServeHTTP(w, r)
	})
}

func GZipWriteHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if client accepts encoding
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Creating gzip.Writer over 'w'
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// Pass gzip-wrapped response writer to next handler
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func GZipReadHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if client sent gzip-encoded request
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Creating gzip.Reader for request body
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
// Also we add userid value to context for handlers to operate
func CheckSessionCookies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Reading cookies
		useridCookie, err1 := r.Cookie("userid")
		shaCookie, err2 := r.Cookie("userid-sha")
		var userid string
		// If we don't have cookie, or we have wrong cookie we have to set it
		if err1 != nil || err2 != nil || !UrlSvc.CheckValidAuth(useridCookie.Value, shaCookie.Value) {
			var err error
			userid, err = UrlSvc.GenerateUserID()
			if err != nil {
				http.Error(w, "error generating user id: "+err.Error(), http.StatusInternalServerError)
				return
			}
			cookie1 := &http.Cookie{
				Name:  "userid",
				Value: userid,
			}

			cookie2 := &http.Cookie{
				Name:  "userid-sha",
				Value: UrlSvc.GenerateSha(userid),
			}
			http.SetCookie(w, cookie1)
			http.SetCookie(w, cookie2)
		} else {
			// If we already have correct cookie we leave userid as received from client
			userid = useridCookie.Value
		}
		ctx := context.WithValue(r.Context(), ContextKeyUserID, userid)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
