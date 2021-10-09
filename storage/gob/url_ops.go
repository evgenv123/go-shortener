package gob

import (
	"context"
	"errors"
	"github.com/evgenv123/go-shortener/model"
	"github.com/evgenv123/go-shortener/storage"
)

func (st *Storage) AddNewURL(ctx context.Context, url model.ShortenedURL) (model.ShortenedURL, error) {
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
	result := model.ShortenedURL{
		ShortURL: shortURLID,
		LongURL:  val.URL,
		UserID:   val.UserID,
	}
	return &result, nil
}

func (st *Storage) GetIDByFull(ctx context.Context, fullURL string) (*model.ShortenedURL, error) {
	for key, val := range st.db.URL {
		if val.URL == fullURL {
			result := model.ShortenedURL{
				ShortURL: key,
				LongURL:  val.URL,
				UserID:   val.UserID,
			}
			return &result, nil
		}
	}
	return nil, errors.New("no item found for full url: " + fullURL)
}

func (st *Storage) GetUserURLs(ctx context.Context, userID string) ([]model.ShortenedURL, error) {
	var result []model.ShortenedURL
	for key, val := range st.db.URL {
		if val.UserID == userID {
			item := model.ShortenedURL{
				ShortURL: key,
				LongURL:  val.URL,
				UserID:   val.UserID,
			}
			result = append(result, item)
		}
	}

	return result, nil
}

func (st *Storage) Ping(ctx context.Context) bool {
	return true
}
