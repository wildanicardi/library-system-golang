package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var mysqlDB *sql.DB

type User struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Image   string `json:"image"`
}
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

// endpoint routes
func Routes() *mux.Router {
	router := mux.NewRouter()
	//user endpoint
	router.HandleFunc("/api/users", indexUser).Methods("GET")
	router.HandleFunc("/api/users", createUser).Methods("POST")
	router.HandleFunc("/api/user/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/api/user/{id}", deleteUser).Methods("DELETE")
	router.HandleFunc("/api/user/{id}", showUser).Methods("GET")
	// books endpoint
	router.HandleFunc("/api/books", indexBook).Methods("GET")
	router.HandleFunc("/api/books", createBook).Methods("POST")
	router.HandleFunc("/api/book/{id}", deleteBook).Methods("DELETE")
	router.HandleFunc("/api/book/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/api/book/{id}", showBook).Methods("GET")

	//transaction
	router.HandleFunc("/api/borrow", indexBorrow).Methods("GET")
	router.HandleFunc("/api/borrow/{idBook}", createBorrow).Methods("POST")
	return router
}

// response json
func renderJSON(res http.ResponseWriter, statusCode int, data interface{}) {
	res.WriteHeader(statusCode)
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(data)
}
func indexBorrow(res http.ResponseWriter, req *http.Request) {

}
func createBorrow(res http.ResponseWriter, req *http.Request) {
	bookID := mux.Vars(req)["idBook"]
	var transaction Transaction
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Print(err)
		return
	}
	datetime := time.Now().Format(time.RFC3339)
	err = json.Unmarshal(body, &transaction)
	if err != nil {
		log.Print(err)
		return
	}
	_, err = mysqlDB.Exec("UPDATE books SET stock = stock - 1 WHERE id = ? AND stock > 0", bookID)
	if err != nil {
		renderJSON(res, http.StatusBadRequest, map[string]interface{}{
			"message": "Failed loan",
		})
	}
	_, err = mysqlDB.Exec("INSERT INTO transaction(user_id,book_id,status,date) VALUES(?,?,?,?)", transaction.User.ID, bookID, transaction.Status, datetime)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	renderJSON(res, http.StatusCreated, map[string]interface{}{
		"message": "Successful loan",
	})
}

// User Function
func indexUser(res http.ResponseWriter, req *http.Request) {
	rows, err := mysqlDB.Query("SELECT id,name,email,address FROM users")
	if err != nil {
		renderJSON(res, http.StatusBadRequest, map[string]interface{}{
			"Message": "Not Found",
		})
	}
	var users []*User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Address); err != nil {
			log.Print(err)
			return
		} else {
			users = append(users, &user)
		}
	}
	renderJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Users",
		"data":    users,
	})
}
func createUser(res http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(4096)
	file, handler, err := req.FormFile("Image")
	name := req.FormValue("Name")
	email := req.FormValue("Email")
	address := req.FormValue("Address")
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
	_, err = mysqlDB.Exec("INSERT INTO users(name,email,address,image) VALUES(?,?,?,?)", name, email, address, handler.Filename)
	if err != nil {
		log.Print(err)
		return
	}
	renderJSON(res, http.StatusCreated, map[string]interface{}{
		"message": "User Created",
	})
}
func updateUser(res http.ResponseWriter, req *http.Request) {
	userID := mux.Vars(req)["id"]
	var user User
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Print(err)
		return
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Print(err)
		return
	}
	query, err := mysqlDB.Prepare("UPDATE users SET name = ?, email = ?,address = ? WHERE id = ?")
	if err != nil {
		log.Print(err)
		return
	}
	query.Exec(user.Name, user.Email, user.Address, userID)

	renderJSON(res, http.StatusCreated, map[string]interface{}{
		"message": "User Updated",
	})

}
func deleteUser(res http.ResponseWriter, req *http.Request) {
	userID := mux.Vars(req)["id"]
	_, err := mysqlDB.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		log.Print(err)
		return
	}
	renderJSON(res, http.StatusAccepted, map[string]interface{}{
		"message": "User Deleted",
	})
}
func showUser(res http.ResponseWriter, req *http.Request) {
	userID := mux.Vars(req)["id"]

	query, err := mysqlDB.Query("SELECT id, name, email,address FROM users WHERE id = " + userID)
	if err != nil {
		log.Print(err)
		return
	}
	var user User
	for query.Next() {
		if err := query.Scan(&user.ID, &user.Name, &user.Email, &user.Address); err != nil {
			log.Print(err)
		}
	}

	renderJSON(res, http.StatusOK, map[string]interface{}{
		"message": "User Show",
		"data":    user,
	})
}

// books function
func indexBook(res http.ResponseWriter, req *http.Request) {
	rows, err := mysqlDB.Query("SELECT id,title,description,image,stock,created_at,updated_at FROM books")
	if err != nil {
		renderJSON(res, http.StatusBadRequest, map[string]interface{}{
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
	renderJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Books",
		"data":    books,
	})
}
func createBook(res http.ResponseWriter, req *http.Request) {
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
	_, err = mysqlDB.Exec("INSERT INTO books(title,description,image,stock,created_at,updated_at) VALUES(?,?,?,?,?,?)", title, description, handler.Filename, stock, datetime, datetime)
	if err != nil {
		log.Print(err)
		return
	}
	renderJSON(res, http.StatusCreated, map[string]interface{}{
		"message": "Book Created",
	})
}
func updateBook(res http.ResponseWriter, req *http.Request) {
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
	_, err = mysqlDB.Exec("UPDATE books SET title = ?, description = ?,image = ?,stock = ?,created_at = ?,updated_at=? WHERE id = ?", title, description, handler.Filename, stock, datetime, datetime, bookID)
	if err != nil {
		log.Print(err)
		return
	}
	renderJSON(res, http.StatusCreated, map[string]interface{}{
		"message": "Book Updated",
	})

}
func deleteBook(res http.ResponseWriter, req *http.Request) {
	bookID := mux.Vars(req)["id"]
	_, err := mysqlDB.Exec("DELETE FROM books WHERE id = ?", bookID)
	if err != nil {
		log.Print(err)
		return
	}
	renderJSON(res, http.StatusAccepted, map[string]interface{}{
		"message": "Book Deleted",
	})
}
func showBook(res http.ResponseWriter, req *http.Request) {
	bookID := mux.Vars(req)["id"]
	query, err := mysqlDB.Query("SELECT id, title, description,image,stock,created_at,updated_at FROM books WHERE id = " + bookID)
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

	renderJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Book Show",
		"data":    book,
	})
}
func main() {
	mysqlDB = Connect()
	defer mysqlDB.Close()
	router := Routes()
	fmt.Println("server started at localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
