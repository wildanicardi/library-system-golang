package controllers

import (
	"encoding/json"
	"fmt"
	"go-project/database"
	"go-project/helper"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type ReportTransaction struct {
	ID          int64     `json:"id"`
	Status      int64     `json:"status"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

//transaction function
func IndexLoan(res http.ResponseWriter, req *http.Request) {
	var reports []*ReportTransaction
	rows, err := database.MysqlDB.Query("SELECT transaction.id,users.name,users.email,books.title,books.description,transaction.status, transaction.date FROM ((transaction INNER JOIN users ON transaction.user_id = users.id)INNER JOIN books ON transaction.book_id = books.id)WHERE transaction.status = 0")
	if err != nil {
		helper.RenderJSON(res, http.StatusBadRequest, map[string]interface{}{
			"Message": "Not Found",
		})
	}
	for rows.Next() {
		var report ReportTransaction
		if err := rows.Scan(&report.ID, &report.Name, &report.Email, &report.Title, &report.Description, &report.Status, &report.Date); err != nil {
			log.Print(err)
			return
		} else {
			reports = append(reports, &report)
		}
	}
	helper.RenderJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Loan",
		"data":    reports,
	})
}
func CreateLoan(res http.ResponseWriter, req *http.Request) {
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
	_, err = database.MysqlDB.Exec("UPDATE books SET stock = stock - 1 WHERE id = ? AND stock > 0", bookID)
	if err != nil {
		helper.RenderJSON(res, http.StatusBadRequest, map[string]interface{}{
			"message": "Failed loan",
		})
	}
	_, err = database.MysqlDB.Exec("INSERT INTO transaction(user_id,book_id,status,date) VALUES(?,?,?,?)", transaction.User.ID, bookID, transaction.Status, datetime)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	helper.RenderJSON(res, http.StatusCreated, map[string]interface{}{
		"message": "Successful Loan",
	})
}
func IndexReturn(res http.ResponseWriter, req *http.Request) {
	var reports []*ReportTransaction
	rows, err := database.MysqlDB.Query("SELECT transaction.id,users.name,users.email,books.title,books.description,transaction.status, transaction.date FROM ((transaction INNER JOIN users ON transaction.user_id = users.id)INNER JOIN books ON transaction.book_id = books.id)WHERE transaction.status = 1")
	if err != nil {
		helper.RenderJSON(res, http.StatusBadRequest, map[string]interface{}{
			"Message": "Not Found",
		})
	}
	for rows.Next() {
		var report ReportTransaction
		if err := rows.Scan(&report.ID, &report.Name, &report.Email, &report.Title, &report.Description, &report.Status, &report.Date); err != nil {
			log.Print(err)
			return
		} else {
			reports = append(reports, &report)
		}
	}
	helper.RenderJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Return",
		"data":    reports,
	})
}
func CreateReturn(res http.ResponseWriter, req *http.Request) {
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
	_, err = database.MysqlDB.Exec("UPDATE books SET stock = stock + 1 WHERE id = ?", bookID)
	if err != nil {
		helper.RenderJSON(res, http.StatusBadRequest, map[string]interface{}{
			"message": "Failed Return",
		})
	}
	_, err = database.MysqlDB.Exec("INSERT INTO transaction(user_id,book_id,status,date) VALUES(?,?,?,?)", transaction.User.ID, bookID, transaction.Status, datetime)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	helper.RenderJSON(res, http.StatusCreated, map[string]interface{}{
		"message": "Successful Return",
	})
}
