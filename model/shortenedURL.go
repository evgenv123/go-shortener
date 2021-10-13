package model

import "time"

type (
	// ShortenedURL represents our canonical data model
	ShortenedURL struct {
		ShortURL  ShortID `json:"short_url"`
		LongURL   string  `json:"original_url"`
		UserID    string
		DeletedAt time.Time
	}
	ShortID int
)
