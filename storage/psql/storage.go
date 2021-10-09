package psql

import (
	"context"
	_ "database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"time"
)

type (
	Storage struct {
		config Config
		db     *sqlx.DB
	}
)

func New(c Config) (*Storage, error) {
	var err error
	st := &Storage{config: c}

	st.db, err = sqlx.Open("pgx", st.config.DSN)
	if err != nil {
		return nil, err
	}

	// Initializing DB table if it does not exist
	// TODO: opt args for ctx!
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = st.db.ExecContext(ctx, initTableCommand)
	if err != nil {
		return nil, fmt.Errorf("error initializing table: %w", err)
	}

	return st, nil
}

func (st Storage) Close() error {
	if st.db == nil {
		return nil
	}

	return st.db.Close()
}

func (st Storage) Ping(ctx context.Context) bool {
	if err := st.db.PingContext(ctx); err != nil {
		return false
	}
	return true
}
