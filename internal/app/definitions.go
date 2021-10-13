package app

import (
	"encoding/json"
	"github.com/evgenv123/go-shortener/internal/config"
	"github.com/evgenv123/go-shortener/service"
	"io"
	"net/http"
)

var appConf config.Config
var URLSvc *service.Processor

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

// InputDelete defines input json format for /api/user/urls endpoint
type InputDelete struct {
	ShortID int
}

func (inp *InputDelete) UnmarshalJSON(data []byte) error {
	var v []interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	inp.ShortID = v[0].(int)

	return nil
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

type ContextKey int

const ContextKeyUserID ContextKey = iota
