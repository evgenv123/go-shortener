package psql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/evgenv123/go-shortener/model"
	"github.com/evgenv123/go-shortener/storage"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

// GetFullByID implements storage.URLReader interface
func (st Storage) GetFullByID(ctx context.Context, shortURLID model.ShortID) (*model.ShortenedURL, error) {
	var fullURL string
	var userID string
	query := "SELECT full_url,user_id FROM " + TableName + " WHERE short_url_id=$1"
	err := st.db.QueryRowContext(ctx, query, int(shortURLID)).Scan(&fullURL, &userID)
	if err != nil {
		// If FULL URL not found return custom error
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.NewFullURLNotFoundErr(shortURLID, err)
		}
		return nil, err
	}
	result := model.ShortenedURL{LongURL: fullURL, ShortURL: shortURLID, UserID: userID}

	return &result, nil
}

// GetIDByFull implements storage.URLReader interface
// Returns sql.ErrNoRows if not found
func (st Storage) GetIDByFull(ctx context.Context, fullURL string) (*model.ShortenedURL, error) {
	var shortID int
	var userID string
	query := "SELECT short_url_id,user_id FROM " + TableName + " WHERE full_url=$1"
	err := st.db.QueryRowContext(ctx, query, fullURL).Scan(&shortID, &userID)

	if err != nil {
		return nil, err
	}
	result := model.ShortenedURL{LongURL: fullURL, ShortURL: model.ShortID(shortID), UserID: userID}

	return &result, nil
}

// GetUserURLs implements storage.URLReader interface
func (st Storage) GetUserURLs(ctx context.Context, userID string) ([]model.ShortenedURL, error) {
	var res []ShortenedURL
	err := st.db.SelectContext(ctx, &res, "SELECT * FROM "+TableName+" WHERE user_id = $1", userID)
	if err != nil {
		// Looks like sqlx.SelectContext does not return error if we have empty result??
		//if errors.Is(err, sql.ErrNoRows) {
		//	return nil, storage.ErrNoURLsForUser
		//}
		return nil, err
	}
	if len(res) == 0 {
		return nil, storage.ErrNoURLsForUser
	}
	// Converting DB output to canonical model
	var ret []model.ShortenedURL
	for _, v := range res {
		ret = append(ret, model.ShortenedURL{
			ShortURL: model.ShortID(v.ShortURL),
			LongURL:  v.LongURL,
			UserID:   v.UserID,
		})
	}
	return ret, nil
}

// AddNewURL implements storage.URLWriter interface
func (st Storage) AddNewURL(ctx context.Context, url model.ShortenedURL) (model.ShortenedURL, error) {
	var result model.ShortenedURL
	query := "INSERT INTO " + TableName + " VALUES ($1, $2, $3) RETURNING *"
	err := st.db.QueryRowContext(ctx, query, int(url.ShortURL), url.LongURL, url.UserID).
		Scan(&result.ShortURL, &result.LongURL, &result.UserID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return result, storage.ErrFullURLExists
			}
		}
		return result, err
	}

	return result, nil
}

// AddBatchURL implements storage.URLWriter interface
func (st Storage) AddBatchURL(ctx context.Context, urls []model.ShortenedURL) error {
	// Starting transaction
	tx, err := st.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Defining statement
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO "+TableName+" VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range urls {
		if _, err = stmt.ExecContext(ctx, int(v.ShortURL), v.LongURL, v.UserID); err != nil {
			return err
		}
	}

	return tx.Commit()
}
