package app

import (
	"math/rand"
	"strconv"
)

// getShortenedURL returns full short URL including server address (using short url id as input)
func getShortenedURL(shortID int) string {
	return appConf.BaseURL + "/" + strconv.Itoa(shortID)
}

func shortenURL(url string) OutputShortURL {
	// Generating ID for link (b)
	var idForLink int
	DB.Lock()
	// Check on duplicate IDs
	for {
		idForLink = rand.Intn(999999)
		_, ok := DB.URLMap[idForLink]
		// If element does not exist we quit
		if !ok {
			break
		}
	}
	DB.URLMap[idForLink] = MappedURL{url, 123}
	DB.Unlock()

	return OutputShortURL{Result: getShortenedURL(idForLink)}
}
