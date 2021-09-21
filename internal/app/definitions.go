package app

import (
	"github.com/evgenv123/go-shortener/internal/config"
	"io"
	"net/http"
	"sync"
)

var appConf config.Config
var DB = ShortenedURLs{URLMap: make(map[int]string)}

type ShortenedURLs struct {
	sync.RWMutex
	URLMap map[int]string
}
type InputURL struct {
	URL string `json:"url"`
}
type OutputShortURL struct {
	Result string `json:"result"`
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}
