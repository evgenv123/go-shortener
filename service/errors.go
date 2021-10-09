package service

import "fmt"

type InvalidURLError struct {
	URL string
	Err error
}

func (myErr *InvalidURLError) Error() string {
	return fmt.Sprintf("%v is not a valid URL %v", myErr.URL, myErr.Err)
}

func (myErr *InvalidURLError) Unwrap() error {
	return myErr.Err
}

func NewInvalidURLError(url string, err error) error {
	return &InvalidURLError{
		URL: url,
		Err: err,
	}
}

type DuplicateFullURLErr struct {
	FullURL  string
	ShortURL string
	Err      error
}

func (myErr *DuplicateFullURLErr) Error() string {
	return fmt.Sprintf("%v already found in storage as %v", myErr.FullURL, myErr.ShortURL)
}

func (myErr *DuplicateFullURLErr) Unwrap() error {
	return myErr.Err
}

func NewDuplicateFullURLErr(fullUrl string, shortUrl string, err error) error {
	return &DuplicateFullURLErr{
		FullURL:  fullUrl,
		ShortURL: shortUrl,
		Err:      err,
	}
}
