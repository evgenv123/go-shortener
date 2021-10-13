package storage

import (
	"errors"
	"fmt"
	"github.com/evgenv123/go-shortener/model"
)

var (
	ErrNoURLsForUser = errors.New("no URLs for user")
	ErrFullURLExists = errors.New("full URL already exists")
	ErrItemDeleted   = errors.New("item deleted")
)

type FullURLNotFoundErr struct {
	ShortURL model.ShortID
	Err      error
}

func (myErr *FullURLNotFoundErr) Error() string {
	return fmt.Sprintf("full url not found for short id %v", int(myErr.ShortURL))
}

func (myErr *FullURLNotFoundErr) Unwrap() error {
	return myErr.Err
}

func NewFullURLNotFoundErr(shortURL model.ShortID, err error) error {
	return &FullURLNotFoundErr{
		ShortURL: shortURL,
		Err:      err,
	}
}
