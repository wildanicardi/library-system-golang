package controllers

import (
	"go-project/database"
	"go-project/helper"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

type Book struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Stock       int64     `json:"stock"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
}
type Transaction struct {
	ID     int64     `json:"id"`
	Status int64     `json:"status"`
	User   *User     `json:"user"`
	Book   *Book     `json:"book"`
	Date   time.Time `json:"date"`
}

func IndexBook(res http.ResponseWriter, req *http.Request) {
	sql := "SELECT id,title,description,image,stock,created_at,updated_at FROM books ORDER BY id DESC"
	rows, err := database.MysqlDB.Query(sql)
	if err != nil {
		helper.RenderJSON(res, http.StatusBadRequest, map[string]interface{}{
			"Message": "Not Found",
		})
	}
	var books []*Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Description, &book.Image, &book.Stock, &book.Created_at, &book.Updated_at); err != nil {
			log.Print(err)
			return
		} else {
			books = append(books, &book)
		}
	}
	helper.RenderJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Books",
		"data":    books,
	})
}
func CreateBook(res http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(4096)
	file, handler, err := req.FormFile("Image")
	title := req.FormValue("Title")
	description := req.FormValue("Description")
	stock := req.FormValue("Stock")
	datetime := time.Now().Format(time.RFC3339)
	if err != nil {
		log.Print(err)
	}
	defer file.Close()
	dir, err := os.Getwd()

	if err != nil {
		log.Print(err)
	}
	fileLocation := filepath.Join(dir, "images", handler.Filename)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Print(err)
		return
	}
	defer targetFile.Close()
	sql := "INSERT INTO books(title,description,image,stock,created_at,updated_at) VALUES(?,?,?,?,?,?)"
	_, err = database.MysqlDB.Exec(sql, title, description, handler.Filename, stock, datetime, datetime)
	if err != nil {
		log.Print(err)
		return
	}
	helper.RenderJSON(res, http.StatusCreated, map[string]interface{}{
		"message": "Book Created",
	})
}
func UpdateBook(res http.ResponseWriter, req *http.Request) {
	bookID := mux.Vars(req)["id"]
	req.ParseMultipartForm(4096)
	file, handler, err := req.FormFile("Image")
	title := req.FormValue("Title")
	description := req.FormValue("Description")
	stock := req.FormValue("Stock")
	datetime := time.Now().Format(time.RFC3339)
	if err != nil {
		log.Print(err)
	}
	defer file.Close()
	dir, err := os.Getwd()

	if err != nil {
		log.Print(err)
	}
	fileLocation := filepath.Join(dir, "images", handler.Filename)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Print(err)
		return
	}
	defer targetFile.Close()
	sql := "UPDATE books SET title = ?, description = ?,image = ?,stock = ?,created_at = ?,updated_at=? WHERE id = ?"
	_, err = database.MysqlDB.Exec(sql, title, description, handler.Filename, stock, datetime, datetime, bookID)
	if err != nil {
		log.Print(err)
		return
	}
	helper.RenderJSON(res, http.StatusCreated, map[string]interface{}{
		"message": "Book Updated",
	})

}
func DeleteBook(res http.ResponseWriter, req *http.Request) {
	bookID := mux.Vars(req)["id"]
	sql := "DELETE FROM books WHERE id = ?"
	_, err := database.MysqlDB.Exec(sql, bookID)
	if err != nil {
		log.Print(err)
		return
	}
	helper.RenderJSON(res, http.StatusAccepted, map[string]interface{}{
		"message": "Book Deleted",
	})
}
func ShowBook(res http.ResponseWriter, req *http.Request) {
	bookID := mux.Vars(req)["id"]
	query, err := database.MysqlDB.Query("SELECT id, title, description,image,stock,created_at,updated_at FROM books WHERE id =  ORDER BY stock DESC" + bookID)
	if err != nil {
		log.Print(err)
		return
	}
	var book Book
	for query.Next() {
		if err := query.Scan(&book.ID, &book.Title, &book.Description, &book.Image, &book.Stock, &book.Created_at, &book.Updated_at); err != nil {
			log.Print(err)
		}
	}

	helper.RenderJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Book Show",
		"data":    book,
	})
}

// stock function
func IndexStock(res http.ResponseWriter, req *http.Request) {
	sql := "SELECT id,title,description,image,stock,created_at,updated_at FROM books"
	rows, err := database.MysqlDB.Query(sql)
	if err != nil {
		helper.RenderJSON(res, http.StatusBadRequest, map[string]interface{}{
			"Message": "Not Found",
		})
	}
	var stocks []*Book
	for rows.Next() {
		var stock Book
		if err := rows.Scan(&stock.ID, &stock.Title, &stock.Description, &stock.Image, &stock.Stock, &stock.Created_at, &stock.Updated_at); err != nil {
			log.Print(err)
			return
		} else {
			stocks = append(stocks, &stock)
		}
	}
	helper.RenderJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Stocks",
		"data":    stocks,
	})
}
