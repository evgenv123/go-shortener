package gob

import (
	"context"
	"errors"
	"github.com/evgenv123/go-shortener/model"
	"github.com/evgenv123/go-shortener/storage"
	"time"
)

func (st *Storage) AddNewURL(ctx context.Context, url model.ShortenedURL) (model.ShortenedURL, error) {
	// Check if full URL already exists
	if _, err := st.GetIDByFull(ctx, url.LongURL); err == nil {
		return model.ShortenedURL{}, storage.ErrFullURLExists
	}
	longURL := LongURL{
		UserID: url.UserID,
		URL:    url.LongURL,
	}
	st.Lock()
	st.db.URL[url.ShortURL] = longURL
	st.Unlock()

	return url, nil
}

func (st *Storage) AddBatchURL(ctx context.Context, urls []model.ShortenedURL) error {
	for _, v := range urls {
		if _, err := st.AddNewURL(ctx, v); err != nil {
			return err
		}
	}

	return nil
}

func (st *Storage) GetFullByID(ctx context.Context, shortURLID model.ShortID) (*model.ShortenedURL, error) {
	st.RLock()
	defer st.RUnlock()
	var val LongURL
	var ok bool
	if val, ok = st.db.URL[shortURLID]; !ok {
		return nil, storage.NewFullURLNotFoundErr(shortURLID, nil)
	}
	// Converting to canonical model
	result, err := val.ToCanonical(shortURLID)

	return &result, err
}

func (st *Storage) GetIDByFull(ctx context.Context, fullURL string) (*model.ShortenedURL, error) {
	for key, val := range st.db.URL {
		if val.URL == fullURL && val.DeletedAt.IsZero() {
			result, err := val.ToCanonical(key)
			return &result, err
		}
	}

	return nil, errors.New("no item found for full url: " + fullURL)
}

func (st *Storage) GetUserURLs(ctx context.Context, userID string) ([]model.ShortenedURL, error) {
	var result []model.ShortenedURL
	for key, val := range st.db.URL {
		if val.UserID == userID && val.DeletedAt.IsZero() {
			item, err := val.ToCanonical(key)
			if err == nil {
				result = append(result, item)
			} else {
				return nil, err
			}
		}
	}
	if len(result) == 0 {
		return result, storage.ErrNoURLsForUser
	}

	return result, nil
}

func (st *Storage) Ping(ctx context.Context) bool {
	return true
}

// DeleteURL implements storage.URLWriter interface
func (st *Storage) DeleteURL(ctx context.Context, url model.ShortenedURL) error {
	st.Lock()
	defer st.Unlock()
	if val, ok := st.db.URL[url.ShortURL]; ok {
		val.DeletedAt = time.Now()
		st.db.URL[url.ShortURL] = val
		return nil
	} else {
		return errors.New("error accessing DB object")
	}
}

// DeleteURL implements storage.URLWriter interface
func (st *Storage) DeleteBatchURL(ctx context.Context, urls []model.ShortenedURL) error {
	for _, v := range urls {
		if err := st.DeleteURL(ctx, v); err != nil {
			return err
		}
	}

	return nil
}
