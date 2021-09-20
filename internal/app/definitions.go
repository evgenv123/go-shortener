package app

import (
	"github.com/evgenv123/go-shortener/internal/config"
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
