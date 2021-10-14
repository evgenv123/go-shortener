package storage

import (
	"context"
	"github.com/evgenv123/go-shortener/model"
	"io"
)

// URLReader defines model.ShortenedURL read operations
type URLReader interface {
	io.Closer
	// GetFullByID finds ShortenedURL object for shortURLID.
	// Returns error not found if nothing found
	GetFullByID(ctx context.Context, shortURLID model.ShortID) (*model.ShortenedURL, error)
	// GetIDByFull finds ID of full URL and returns corresponding model.ShortenedURL object.
	// Returns error not found if nothing found
	GetIDByFull(ctx context.Context, fullURL string) (*model.ShortenedURL, error)
	// GetUserURLs finds all ShortenedURLs associated with userID
	// If no URLs found returns storage.NoURLsForUserErr error
	GetUserURLs(ctx context.Context, userID string) ([]model.ShortenedURL, error)
	// Ping checks if the storage is online
	Ping(ctx context.Context) bool
}

// URLWriter defines model.ShortenedURL write operations
type URLWriter interface {
	io.Closer
	// AddNewURL adds new ShortenedURL
	AddNewURL(ctx context.Context, url model.ShortenedURL) (model.ShortenedURL, error)
	// AddBatchURL adds array of new ShortenedURLs
	AddBatchURL(ctx context.Context, urls []model.ShortenedURL) error
	// Ping checks if the storage is online
	Ping(ctx context.Context) bool
	DeleteURL(ctx context.Context, url model.ShortenedURL) error
	DeleteBatchURL(ctx context.Context, urls []model.ShortenedURL) error
}
