package gob

import "github.com/evgenv123/go-shortener/model"

type (
	// ShortenedURLs represents model.ShortenedURL canonical model for GOB storage
	ShortenedURLs struct {
		URL map[model.ShortID]LongURL
	}

	LongURL struct {
		URL    string
		UserID string
	}
)
