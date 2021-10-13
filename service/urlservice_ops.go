package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/evgenv123/go-shortener/model"
	"github.com/evgenv123/go-shortener/storage"
	"github.com/google/uuid"
	"math/rand"
	"net/url"
	"strconv"
)

// CheckValidAuth checks if 'sha' is correct signature for 'userid'
func (svc *Processor) CheckValidAuth(userid string, sha string) bool {
	return svc.GenerateSha(userid) == sha
}

func (svc *Processor) GetUserURLs(ctx context.Context, userID string) ([]model.ShortenedURL, error) {
	result, err := svc.urlStorage.GetUserURLs(ctx, userID)
	// Switching error type from storage to service
	// no need for this
	//if errors.Is(err, storage.ErrNoURLsForUser) {
	//	return result, ErrNoURLsForUser
	//}
	return result, err
}

func (svc *Processor) Ping(ctx context.Context) bool {
	return svc.urlStorage.Ping(ctx)
}

// GenerateSha generates sha256 base64-encoded hash for userid
func (svc *Processor) GenerateSha(userid string) string {
	// Skipping error check because we did it in config.Validate()
	secret, _ := hex.DecodeString(svc.config.HexSecret)
	// Calculating valid hash
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(userid))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// GenerateUserID generates uuid
func (svc *Processor) GenerateUserID() (string, error) {
	useruuid, err := uuid.NewRandom()
	return useruuid.String(), err
}

// generateNewID generates new unique id for short URL
func (svc *Processor) generateNewID(ctx context.Context) (model.ShortID, error) {
	var idForLink model.ShortID
	for {
		idForLink = model.ShortID(rand.Intn(999999))
		// If ID is not duplicate (if we didn't find full url in storage)
		var myErr *storage.FullURLNotFoundErr
		if _, err := svc.urlStorage.GetFullByID(ctx, idForLink); errors.As(err, &myErr) {
			break
		} else if err != nil {
			return 0, err
		}
	}
	return idForLink, nil
}

// GetFullLinkShortObj returns full URI of short link from short url object
func (svc *Processor) GetFullLinkShortObj(ctx context.Context, shortObj *model.ShortenedURL) string {
	return svc.config.BaseURL + "/" + strconv.Itoa(int(shortObj.ShortURL))
}

// GetObjFromShortID finds model.ShortenedURL object for corresponding short link id
func (svc *Processor) GetObjFromShortID(ctx context.Context, shortID model.ShortID) (*model.ShortenedURL, error) {
	result, err := svc.urlStorage.GetFullByID(ctx, shortID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetObjFromFullURL return error not found, object or other error
func (svc *Processor) GetObjFromFullURL(ctx context.Context, fullURL string) (*model.ShortenedURL, error) {
	return svc.urlStorage.GetIDByFull(ctx, fullURL)
}

func (svc *Processor) ShortenURL(ctx context.Context, fullURL string, userID string) (*model.ShortenedURL, error) {
	// Checking for valid URL
	_, err := url.ParseRequestURI(fullURL)
	if err != nil {
		return nil, NewInvalidURLError(fullURL, err)
	}
	idForLink, err := svc.generateNewID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error generating short ID: %w", err)
	}
	// Trying to add URL
	result, err := svc.urlStorage.AddNewURL(ctx, model.ShortenedURL{ShortURL: idForLink, LongURL: fullURL, UserID: userID})
	// If we get item already exists error we send back error and existing item
	if errors.Is(err, storage.ErrFullURLExists) {
		obj, err := svc.GetObjFromFullURL(ctx, fullURL)
		return obj, NewDuplicateFullURLErr(fullURL, svc.GetFullLinkShortObj(ctx, obj), err)
	}

	return &result, err
}

// DeleteBatchURL asynchronously deletes urls from DB
func (svc *Processor) DeleteBatchURL(ctx context.Context, urls []model.ShortenedURL) error {
	for _, v := range urls {
		if err := svc.urlStorage.DeleteURL(ctx, v); err != nil {
			return err
		}
	}
	return nil
}
