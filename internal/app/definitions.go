package app

import (
	"github.com/evgenv123/go-shortener/internal/config"
	"github.com/evgenv123/go-shortener/service"
	"io"
	"net/http"
)

var appConf config.Config
var UrlSvc *service.Processor

// InputURL defines json input format for /api/shorten endpoint
type InputURL struct {
	URL string `json:"url"`
}
type OutputShortURL struct {
	Result string `json:"result"`
}

// OutputAllURLs defines json output format for /user/urls handler
type OutputAllURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// InputBatch defines input json format for /api/shorten/batch endpoint
type InputBatch struct {
	CorrID  string `json:"correlation_id"`
	OrigURL string `json:"original_url"`
}

// OutputBatch defines output json format for /api/shorten/batch endpoint
type OutputBatch struct {
	CorrID   string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

type ContextKey int

const ContextKeyUserID ContextKey = iota
