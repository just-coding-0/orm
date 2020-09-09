package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)



func Open(url string) (*sql.DB, error) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	return db, nil
}
