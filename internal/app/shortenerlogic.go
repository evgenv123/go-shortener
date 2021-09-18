package app

import (
	"math/rand"
	"strconv"
)

func shortenURL(url string) OutputShortURL {
	// Generating ID for link (b)
	var idForLink int
	DB.mu.Lock()
	// Check on duplicate IDs
	for {
		idForLink = rand.Intn(999999)
		_, ok := DB.URLMap[idForLink]
		// If element does not exist we quit
		if !ok {
			break
		}
	}
	DB.URLMap[idForLink] = url
	DB.mu.Unlock()

	return OutputShortURL{Result: appConf.BaseURL + "/" + strconv.Itoa(idForLink)}
}
