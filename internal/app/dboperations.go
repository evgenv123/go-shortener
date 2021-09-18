package app

import (
	"encoding/gob"
	"os"
)

func WriteDBToFile() error {
	// create a file
	dataFile, err := os.Create(appConf.FileStorage)
	if err != nil {
		return err
	}
	defer dataFile.Close()
	// serialize the data
	dataEncoder := gob.NewEncoder(dataFile)
	DB.mu.RLock()
	err = dataEncoder.Encode(DB)
	DB.mu.RUnlock()
	if err != nil {
		return err
	}
	return nil
}

// readDBFromFile return nil is file does not exist or read successfully
// Return error if cannot read file
func readDBFromFile() error {
	// open data file
	dataFile, err := os.Open(appConf.FileStorage)
	// If file does not exist there is no error
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer dataFile.Close()

	dataDecoder := gob.NewDecoder(dataFile)
	DB.mu.Lock()
	err = dataDecoder.Decode(&DB)
	DB.mu.Unlock()
	if err != nil {
		return err
	}
	return nil
}
