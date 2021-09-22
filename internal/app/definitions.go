package app

import (
	"github.com/evgenv123/go-shortener/internal/config"
	"io"
	"net/http"
	"sync"
)

var appConf config.Config
var DB = ShortenedURLs{URLMap: make(map[int]MappedURL)}

type MappedURL struct {
	URL    string
	UserID int
}

type ShortenedURLs struct {
	sync.RWMutex
	URLMap map[int]MappedURL
}
type InputURL struct {
	URL string `json:"url"`
}
type OutputShortURL struct {
	Result string `json:"result"`
}
type OutputAllURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

type gzipReader struct {
	http.Request
	Reader io.Reader
}

type Middleware func(http.Handler) http.Handler
