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
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
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
	return router
}

// response json
func renderJSON(res http.ResponseWriter, data interface{}) {
	res.WriteHeader(200)
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(data)
}

// User Function
func indexUser(res http.ResponseWriter, req *http.Request) {
	rows, err := mysqlDB.Query("SELECT id,name,email,address FROM users")
	if err != nil {
		renderJSON(res, map[string]interface{}{
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
	renderJSON(res, map[string]interface{}{
		"status":  200,
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
	renderJSON(res, map[string]interface{}{
		"status":  201,
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

	renderJSON(res, map[string]interface{}{
		"status":  201,
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
	renderJSON(res, map[string]interface{}{
		"status":  204,
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

	renderJSON(res, map[string]interface{}{
		"status":  200,
		"message": "User Show",
		"data":    user,
	})
}

// books function
func indexBook(res http.ResponseWriter, req *http.Request) {
	rows, err := mysqlDB.Query("SELECT id,title,description,image FROM books")
	if err != nil {
		renderJSON(res, map[string]interface{}{
			"Message": "Not Found",
		})
	}
	var books []*Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Description, &book.Image); err != nil {
			log.Print(err)
			return
		} else {
			books = append(books, &book)
		}
	}
	renderJSON(res, map[string]interface{}{
		"status":  200,
		"message": "Books",
		"data":    books,
	})
}
func createBook(res http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(4096)
	file, handler, err := req.FormFile("Image")
	title := req.FormValue("Title")
	description := req.FormValue("Description")
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
	_, err = mysqlDB.Exec("INSERT INTO books(title,description,image) VALUES(?,?,?)", title, description, handler.Filename)
	if err != nil {
		log.Print(err)
		return
	}
	renderJSON(res, map[string]interface{}{
		"status":  201,
		"message": "Book Created",
	})
}
func updateBook(res http.ResponseWriter, req *http.Request) {
	bookID := mux.Vars(req)["id"]
	req.ParseMultipartForm(4096)
	file, handler, err := req.FormFile("Image")
	title := req.FormValue("Title")
	description := req.FormValue("Description")
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
	_, err = mysqlDB.Exec("UPDATE books SET title = ?, description = ?,image = ? WHERE id = ?", title, description, handler.Filename, bookID)
	if err != nil {
		log.Print(err)
		return
	}
	renderJSON(res, map[string]interface{}{
		"status":  201,
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
	renderJSON(res, map[string]interface{}{
		"status":  204,
		"message": "Book Deleted",
	})
}
func showBook(res http.ResponseWriter, req *http.Request) {
	bookID := mux.Vars(req)["id"]

	query, err := mysqlDB.Query("SELECT id, title, description,image FROM books WHERE id = " + bookID)
	if err != nil {
		log.Print(err)
		return
	}
	var book Book
	for query.Next() {
		if err := query.Scan(&book.ID, &book.Title, &book.Description, &book.Image); err != nil {
			log.Print(err)
		}
	}

	renderJSON(res, map[string]interface{}{
		"status":  200,
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
