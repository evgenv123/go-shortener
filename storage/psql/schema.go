package psql

import (
	"database/sql"
	"github.com/evgenv123/go-shortener/model"
)

const (
	TableName        = "shortURLs"
	initTableCommand = `
-- Short URLs table
create table if not exists ` + TableName + `
(
	short_url_id	int,
    full_url		varchar(255) not null,
    user_id			varchar(100) not null,
	deleted_at		timestamp with time zone,
    unique (short_url_id),
	unique (full_url)
);
`
)

type (
	// ShortenedURL represents model.ShortenedURL canonical model for PSQL storage
	ShortenedURL struct {
		ShortURL  int          `db:"short_url_id"`
		LongURL   string       `db:"full_url"`
		UserID    string       `db:"user_id"`
		DeletedAt sql.NullTime `db:"deleted_at"`
	}

	ShortenedURLs []ShortenedURL
)

// ToCanonical converts ShortenedURLs to canonical model []model.ShortenedURL
func (u ShortenedURLs) ToCanonical() ([]model.ShortenedURL, error) {
	var ret []model.ShortenedURL
	for _, v := range u {
		if newItem, err := v.ToCanonical(); err == nil {
			ret = append(ret, newItem)
		} else {
			return nil, err
		}
	}
	return ret, nil
}

// ToCanonical converts ShortenedURL to canonical model model.ShortenedURL
func (u ShortenedURL) ToCanonical() (model.ShortenedURL, error) {
	ret := model.ShortenedURL{
		ShortURL:  model.ShortID(u.ShortURL),
		LongURL:   u.LongURL,
		UserID:    u.UserID,
		DeletedAt: u.DeletedAt.Time,
	}

	return ret, nil
}
