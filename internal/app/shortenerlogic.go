package app

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/evgenv123/go-shortener/internal/dbcore"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"log"
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

// GetShortenedURL returns full short URL including server address (using short url id as input)
func GetShortenedURL(shortID int) string {
	return appConf.BaseURL + "/" + strconv.Itoa(shortID)
}

func shortenURL(url string, userid string) (string, error) {
	// Generating ID for link (b)
	var idForLink int
	DB.RLock()
	// Check on duplicate IDs
	for {
		idForLink = rand.Intn(999999)
		_, ok := DB.URLMap[idForLink]
		// If element does not exist we quit
		if !ok {
			break
		}
	}
	DB.RUnlock()

	if err := dbcore.InsertURL(url, idForLink, userid); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				// Here we request existing short url for our full url
				short, err2 := dbcore.GetShortFromFull(url)
				err = NewFullURLDuplicateError(url, GetShortenedURL(short), err)
				// If we could not retrieve short url id for full url we wrap this error too
				if err2 != nil {
					err = fmt.Errorf("%s %w", err2, err)
				}
			}
		}
		log.Println("Error inserting short URL to DB ", err)
		return "", err
	}

	// Writing new link to file DB
	DB.Lock()
	DB.URLMap[idForLink] = MappedURL{url, userid}
	DB.Unlock()

	return GetShortenedURL(idForLink), nil
}
