package app

import (
	"github.com/evgenv123/go-shortener/internal/config"
	"io"
	"net/http"
	"sync"
)

var appConf config.Config
var DB = ShortenedURLs{URLMap: make(map[int]MappedURL)}

var myDirtyLittleSecret = []byte{18, 232, 139, 12, 216, 189, 22, 128, 122, 49, 246, 137, 191, 24, 38, 210}

type MappedURL struct {
	URL    string
	UserID string
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

// Types for /api/shorten/batch endpoint
type InputBatch struct {
	CorrID  string `json:"correlation_id"`
	OrigURL string `json:"original_url"`
}
type OutputBatch struct {
	CorrID   string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

type Middleware func(http.Handler) http.Handler

type contextKey int

const contextKeyUserID contextKey = iota
