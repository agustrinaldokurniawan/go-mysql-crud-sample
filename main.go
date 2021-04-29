package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type User struct {
	Id    int
	Name  string
	Email string
}

func allUsers(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "All Users Endpoint Hit")

	db := dbConnection()
	allUser, err := db.Query("SELECT * FROM test")
	if err != nil {
		panic(err.Error())
	}

	users := []User{}
	for allUser.Next() {
		var r User
		err := allUser.Scan(&r.Id, &r.Name, &r.Email)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, r)
	}

	json.NewEncoder(w).Encode(users)
	defer db.Close()
}

func newUser(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "New Users Endpoint Hit")

	db := dbConnection()

	vars := mux.Vars(r)
	name := vars["name"]
	email := vars["email"]

	insForm, err := db.Prepare("INSERT INTO test (name, email) VALUES (?,?)")
	if err != nil {
		panic(err.Error())
	}

	insForm.Exec(name, email)

	fmt.Fprintf(w, "New User Successfully Created")

	defer db.Close()
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Delete Users Endpoint Hit")

	db := dbConnection()

	vars := mux.Vars(r)
	id := vars["id"]

	deleteRow, err := db.Prepare("DELETE FROM test WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	deleteRow.Exec(id)

	fmt.Fprintf(w, "User Successfully Deleted")
	defer db.Close()
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Update Users Endpoint Hit")

	db := dbConnection()

	vars := mux.Vars(r)
	email := vars["email"]
	id := vars["id"]

	updateForm, err := db.Prepare("UPDATE test SET email=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	updateForm.Exec(email, id)

	fmt.Fprintf(w, "User Successfully Updated")
	defer db.Close()
}

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/users", allUsers).Methods("GET")
	myRouter.HandleFunc("/user/{id}", deleteUser).Methods("DELETE")
	myRouter.HandleFunc("/user/{id}/{email}", updateUser).Methods("PUT")
	myRouter.HandleFunc("/user/{name}/{email}", newUser).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func dbConnection() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:passwordroot@tcp(127.0.0.1:3306)/test")
	if err != nil {
		fmt.Println(err.Error())
		panic("[mysql] failed to connect database")
	}

	return db
}

func main() {
	fmt.Println("Go Mysql CRUD Sample")

	handleRequest()
}
