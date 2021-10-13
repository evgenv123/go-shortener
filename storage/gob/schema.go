package gob

import (
	"github.com/evgenv123/go-shortener/model"
	"time"
)

type (
	// ShortenedURLs represents model.ShortenedURL canonical model for GOB storage
	ShortenedURLs struct {
		URL map[model.ShortID]LongURL
	}

	LongURL struct {
		URL       string
		UserID    string
		DeletedAt time.Time
	}
)

// ToCanonical converts LongURL to canonical model.ShortenedURL using parameter model.ShortID
func (long LongURL) ToCanonical(short model.ShortID) (model.ShortenedURL, error) {
	// Converting to canonical model
	result := model.ShortenedURL{
		ShortURL:  short,
		LongURL:   long.URL,
		UserID:    long.UserID,
		DeletedAt: long.DeletedAt,
	}

	return result, nil
}
