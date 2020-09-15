package orm

import "database/sql"

type db struct {
	db *sql.DB
}



const (
	Mysql = "mysql"
)

func NewDb(dbType string)  {

}
