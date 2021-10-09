package storage

import (
	"fmt"
	"github.com/evgenv123/go-shortener/model"
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

func NewFullURLNotFoundErr(shortUrl model.ShortID, err error) error {
	return &FullURLNotFoundErr{
		ShortURL: shortUrl,
		Err:      err,
	}
}