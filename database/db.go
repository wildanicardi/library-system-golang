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

	connection := fmt.Sprintf("%s:@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", user, host, port, database)
	db, err := sql.Open("mysql", connection)
	if err != nil {
		log.Print(err, "\nError connect database")
	}
	if db == nil {
		panic("db nil")
	}
	migrate(db)
	return db
}

// migrate db
func migrate(db *sql.DB) {
	sql := `
	CREATE TABLE IF NOT EXISTS books(
		id INTEGER PRIMARY KEY AUTO_INCREMENT,
		title VARCHAR(100) NOT NULL,
		description TEXT NOT NULL,
		image TEXT NOT NULL,
		stock INTEGER(11) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT current_timestamp(),
		updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp()
	);
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) NOT NULL,
		address VARCHAR(100) NOT NULL,
		image TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS transaction(
		id INTEGER PRIMARY KEY AUTO_INCREMENT,
		user_id INTEGER(11) DEFAULT NULL,
		book_id INTEGER(11) DEFAULT NULL,
		status INTEGER(11) NOT NULL,
		date TIMESTAMP NOT NULL DEFAULT current_timestamp(),
		CONSTRAINT FK_UsersTransaction FOREIGN KEY (user_id) REFERENCES Users(id),
		CONSTRAINT FK_BooksTransaction FOREIGN KEY (book_id) REFERENCES Books(id)
	);
	`
	_, err := db.Exec(sql)
	if err != nil {
		log.Print(err)
		return
	}
}
