package database

import (
	"database/sql"
	"fmt"
	"log"
)
var MysqlDB *sql.DB
// connection database mysql
func Connect() *sql.DB {
	user := "root"
	host := "localhost"
	port := "3306"
	database := "library"

	connection := fmt.Sprintf("%s:@tcp(%s:%s)/%s?parseTime=true", user, host, port, database)
	db, err := sql.Open("mysql", connection)
	if err != nil {
		log.Print(err, "\nError connect database")
	}
	return db
}
