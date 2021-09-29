package dbcore

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"time"
)

var db *sql.DB

func Init(dsn string) error {
	var err error
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

func InsertURL(fullUrl string, shortURLID int, userid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err := db.ExecContext(ctx, "INSERT INTO "+TableName+" VALUES (?, ?, ?)", shortURLID, fullUrl, userid)
	if err != nil {
		return err
	}
	return nil
}

func UnshortURL(shortID int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var fullURL string
	err := db.QueryRowContext(ctx, "SELECT full_url FROM "+TableName+" WHERE short_url_id=?", shortID).Scan(&fullURL)

	if err != nil {
		return "", err
	}
	return fullURL, nil
}
