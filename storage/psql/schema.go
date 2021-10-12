package psql

const (
	TableName        = "shortURLs"
	initTableCommand = `
-- Short URLs table
create table if not exists ` + TableName + `
(
	short_url_id	int,
    full_url		varchar(255) not null,
    user_id			varchar(100) not null,
    unique (short_url_id),
	unique (full_url)
);
`
)

type (
	// ShortenedURL represents model.ShortenedURL canonical model for PSQL storage
	ShortenedURL struct {
		ShortURL int    `db:"short_url_id"`
		LongURL  string `db:"full_url"`
		UserID   string `db:"user_id"`
	}
)
