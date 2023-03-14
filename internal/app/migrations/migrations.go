package migrations

import (
	"database/sql"
	"github.com/lopezator/migrator"
)

func Up(db *sql.DB) error {
	m, err := migrator.New(
		migrator.Migrations(
			&migrator.MigrationNoTx{
				Name: "Create urls table",
				Func: createUrlsTable,
			},
			&migrator.MigrationNoTx{
				Name: "Add user_id index to urls table",
				Func: addUserIDIndexToUrlsTable,
			},
			&migrator.MigrationNoTx{
				Name: "Add url unique index to urls table",
				Func: addURLUniqueIndexToUrlsTable,
			},
			&migrator.MigrationNoTx{
				Name: "Add deleted column to urls table",
				Func: addDeletedColumnToUrlsTable,
			},
		),
	)
	if err != nil {
		return err
	}

	return m.Migrate(db)
}

func createUrlsTable(db *sql.DB) error {
	_, err := db.Exec(`
create table urls
(
    id      serial        not null primary key,
    user_id uuid          not null,
    url_id  varchar(16)   not null unique,
    url     varchar(2000) not null
)
	`)

	return err
}

func addUserIDIndexToUrlsTable(db *sql.DB) error {
	_, err := db.Exec("create index urls_user_id_index on urls (user_id)")

	return err
}

func addURLUniqueIndexToUrlsTable(db *sql.DB) error {
	_, err := db.Exec("alter table urls add unique (url)")

	return err
}

func addDeletedColumnToUrlsTable(db *sql.DB) error {
	_, err := db.Exec("alter table urls add deleted bool default false not null")

	return err
}
