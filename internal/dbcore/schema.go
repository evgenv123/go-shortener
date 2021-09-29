package dbcore

const TableName = "shortURLs"

const initTableCommand = `
-- Short URLs table
create table ` + TableName + `
(
	short_url_id	int,
    full_url		varchar(255) not null,
    user_id			varchar(100) not null,
    unique (url),
	unique (user_id)
);
`
