package dbcore

const TableName = "shortURLs"

const initTableCommand = `
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
