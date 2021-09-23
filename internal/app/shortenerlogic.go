package app

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"strconv"
)

func generateSha(userid string) string {
	// Calculating valid hash
	h := hmac.New(sha256.New, myDirtyLittleSecret)
	h.Write([]byte(userid))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// checkValidAuth checks if 'sha' is correct signature for 'userid'
func checkValidAuth(userid string, sha string) bool {
	return generateSha(userid) == sha
}

// getShortenedURL returns full short URL including server address (using short url id as input)
func getShortenedURL(shortID int) string {
	return appConf.BaseURL + "/" + strconv.Itoa(shortID)
}

func shortenURL(url string, userid string) OutputShortURL {
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
	DB.URLMap[idForLink] = MappedURL{url, userid}
	DB.Unlock()

	return OutputShortURL{Result: getShortenedURL(idForLink)}
}
