package model

import (
	"encoding/json"
	"strconv"
	"time"
)

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

// UnmarshalJSON is a custom unmarshal for ShortID type
func (inp *ShortID) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch t := v.(type) {
	case int:
		*inp = ShortID(t)
	case string:
		i, err := strconv.Atoi(t)
		if err != nil {
			return err
		}
		*inp = ShortID(i)
	}

	return nil
}
