package app

import "fmt"

type FullURLDuplicateError struct {
	FullURL  string
	ShortURL string
	Err      error
}

func (myErr *FullURLDuplicateError) Error() string {
	return fmt.Sprintf("%v already has a short link %v %v", myErr.FullURL, myErr.ShortURL, myErr.Err)
}

func NewFullURLDuplicateError(full string, short string, err error) error {
	return &FullURLDuplicateError{
		FullURL:  full,
		ShortURL: short,
		Err:      err,
	}
}
