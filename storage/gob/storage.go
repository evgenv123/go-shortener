package gob

import (
	"encoding/gob"
	"github.com/evgenv123/go-shortener/model"
	"os"
	"sync"
)

type (
	// Storage defines data structure for GOB
	Storage struct {
		sync.RWMutex
		config Config
		db     ShortenedURLs
	}
)

// New creates a new Storage using !validated! Config
func New(c Config) (*Storage, error) {
	st := &Storage{
		config: c,
		db:     ShortenedURLs{URL: make(map[model.ShortID]LongURL)},
	}

	// Trying to read DB from file if exists
	dataFile, err := os.Open(st.config.Filename)
	// If file does not exist there is no error
	if os.IsNotExist(err) {
		return st, nil
	}
	if err != nil {
		return nil, err
	}
	defer dataFile.Close()

	dataDecoder := gob.NewDecoder(dataFile)
	st.Lock()
	err = dataDecoder.Decode(&st.db)
	st.Unlock()
	if err != nil {
		return nil, err
	}

	return st, nil
}

// Close writes DB to file
func (st *Storage) Close() error {
	// create (rewrite) file
	dataFile, err := os.Create(st.config.Filename)
	if err != nil {
		return err
	}
	defer dataFile.Close()
	// serialize the data
	dataEncoder := gob.NewEncoder(dataFile)
	st.RLock()
	err = dataEncoder.Encode(st.db)
	st.RUnlock()
	if err != nil {
		return err
	}

	return nil
}
