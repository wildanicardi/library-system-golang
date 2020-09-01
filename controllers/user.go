package controllers

import (
	"encoding/json"
	"fmt"
	"go-project/database"
	"go-project/helper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

type User struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Image   string `json:"image"`
}

func IndexUser(res http.ResponseWriter, req *http.Request) {
	rows, err := database.MysqlDB.Query("SELECT id,name,email,address,image FROM users")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var users []*User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Address, &user.Image); err != nil {
			fmt.Println(err.Error())
			return
		} else {
			users = append(users, &user)
		}
	}
	helper.RenderJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Users",
		"data":    users,
	})
}
func CreateUser(res http.ResponseWriter, req *http.Request) {
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
	_, err = database.MysqlDB.Exec("INSERT INTO users(name,email,address,image) VALUES(?,?,?,?)", name, email, address, handler.Filename)
	if err != nil {
		log.Print(err)
		return
	}
	helper.RenderJSON(res, http.StatusCreated, map[string]interface{}{
		"message": "User Created",
	})
}
func ShowUser(res http.ResponseWriter, req *http.Request) {
	userID := mux.Vars(req)["id"]

	query, err := database.MysqlDB.Query("SELECT id, name, email,address,image FROM users WHERE id = " + userID)
	if err != nil {
		log.Print(err)
		return
	}
	var user User
	for query.Next() {
		if err := query.Scan(&user.ID, &user.Name, &user.Email, &user.Address, &user.Image); err != nil {
			log.Print(err)
		}
	}

	helper.RenderJSON(res, http.StatusOK, map[string]interface{}{
		"message": "User Show",
		"data":    user,
	})
}
func UpdateUser(res http.ResponseWriter, req *http.Request) {
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
	query, err := database.MysqlDB.Prepare("UPDATE users SET name = ?, email = ?,address = ? WHERE id = ?")
	if err != nil {
		log.Print(err)
		return
	}
	query.Exec(user.Name, user.Email, user.Address, userID)

	helper.RenderJSON(res, http.StatusCreated, map[string]interface{}{
		"message": "User Updated",
	})

}
func DeleteUser(res http.ResponseWriter, req *http.Request) {
	userID := mux.Vars(req)["id"]
	_, err := database.MysqlDB.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		log.Print(err)
		return
	}
	helper.RenderJSON(res, http.StatusAccepted, map[string]interface{}{
		"message": "User Deleted",
	})
}
