package dbcore

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v4/stdlib"
	"time"
)

var db *sql.DB
var dbOnline = false

func Init(dsn string) error {
	var err error
	dbOnline = false
	// dsnExample := "postgres://postgres:tttest@localhost:5432/postgres"
	db, err = sql.Open("pgx", dsn)
	if err != nil {
		return err
	}

	// Initializing DB table if it does not exist
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx, initTableCommand)
	if err != nil {
		return err
	}
	dbOnline = true
	return nil
}

func Close() {
	db.Close()
}

// CheckConn pings the database and returns true if everything is ok and false otherwise
func CheckConn() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return false
	}
	return true
}

func InsertURL(fullURL string, shortURLID int, userid string) error {
	if !dbOnline {
		return errors.New("database offline")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := "INSERT INTO " + TableName + " VALUES ($1, $2, $3)"
	_, err := db.ExecContext(ctx, query, shortURLID, fullURL, userid)
	if err != nil {
		return err
	}
	return nil
}

func UnshortenURL(shortID int) (string, error) {
	if !dbOnline {
		return "", errors.New("database offline")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var fullURL string
	query := "SELECT full_url FROM " + TableName + " WHERE short_url_id=?"
	err := db.QueryRowContext(ctx, query, shortID).Scan(&fullURL)

	if err != nil {
		return "", err
	}
	return fullURL, nil
}

// GetShortFromFull returns short URL ID from DB if it was shortened before
func GetShortFromFull(fullURL string) (int, error) {
	if !dbOnline {
		return 0, errors.New("database offline")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var shortID int
	query := "SELECT short_url_id FROM " + TableName + " WHERE full_url=$1"
	err := db.QueryRowContext(ctx, query, fullURL).Scan(&shortID)

	if err != nil {
		return 0, err
	}
	return shortID, nil
}
