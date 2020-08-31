package main

import (
	"fmt"
	"go-project/controllers"
	"go-project/database"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// var database.MysqlDB *sql.DB

// endpoint routes
func Routes() *mux.Router {
	router := mux.NewRouter()
	//user endpoint
	router.HandleFunc("/api/users", controllers.IndexUser).Methods("GET")
	router.HandleFunc("/api/users", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/api/user/{id}", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/api/user/{id}", controllers.DeleteUser).Methods("DELETE")
	router.HandleFunc("/api/user/{id}", controllers.ShowUser).Methods("GET")
	// books endpoint
	router.HandleFunc("/api/books", controllers.IndexBook).Methods("GET")
	router.HandleFunc("/api/books", controllers.CreateBook).Methods("POST")
	router.HandleFunc("/api/book/{id}", controllers.DeleteBook).Methods("DELETE")
	router.HandleFunc("/api/book/{id}", controllers.UpdateBook).Methods("PUT")
	router.HandleFunc("/api/book/{id}", controllers.ShowBook).Methods("GET")

	//transaction
	router.HandleFunc("/api/loan", controllers.IndexLoan).Methods("GET")
	router.HandleFunc("/api/loan/{idBook}", controllers.CreateLoan).Methods("POST")
	router.HandleFunc("/api/return", controllers.IndexReturn).Methods("GET")
	router.HandleFunc("/api/return/{idBook}", controllers.CreateReturn).Methods("POST")

	//stock
	router.HandleFunc("/api/stocks", controllers.IndexStock).Methods("GET")
	return router
}

func main() {
	database.MysqlDB = database.Connect()
	defer database.MysqlDB.Close()
	router := Routes()
	fmt.Println("server started at localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
