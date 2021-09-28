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
	return err
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
