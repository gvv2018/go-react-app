package main

import (
	"log"
	"net/http"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	_"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Hotdog struct{
	Name string `json:"name"`
	Price string `json:"price"`
}

//var hotdogs []Hotdog

var db *sql.DB
var err error

// Function gets all rows from db

func getHotdogs( w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var hotdogs []Hotdog

	result, err := db.Query("SELECT name, price FROM type_hotdogs;")
	if err != nil {
		panic(err.Error())
	}

//	defer db.Close()

	for result.Next() {
		var hotdog Hotdog
		err := result.Scan(&hotdog.Name, &hotdog.Price)
		if err != nil {
		panic(err.Error())
		}
		hotdogs = append(hotdogs, hotdog)
	}
	json.NewEncoder(w).Encode(hotdogs)

}

// Function adds new rows to db

func createHotdog( w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stmt, err := db.Prepare("INSERT INTO type_hotdogs (name, price) VALUES(?, ?)")
	if err != nil {
	    panic(err.Error())
	}

	body, err := ioutil.ReadAll (r.Body)
	if err != nil {
	    panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	name := keyVal["name"]
	price := keyVal["price"]
	
	_, err = stmt.Exec(name, price)
	if err != nil {
	    panic(err.Error())
	}

//	fmt.Println(name, price)
}

// Function delete row with db

func deleteHotdog( w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	stmt, err := db.Prepare("DELETE FROM type_hotdogs WHERE name = ?")
	if err != nil {
	    panic(err.Error())
	}

	_, err = stmt.Exec(params["name"])
	if err != nil {
	    panic(err.Error())
	}

}

// Function update row oin database

func updateHotdog( w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	stmt, err := db.Prepare("UPDATE type_hotdogs SET name = ?, price = ? WHERE name = ?")
	if err != nil {
	    panic(err.Error())
	}

	body, err := ioutil.ReadAll (r.Body)
	if err != nil {
	    panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newName := keyVal["name"]
	newPrice := keyVal["price"]
	
	_, err = stmt.Exec(newName, newPrice, params["name"])
	if err != nil {
	    panic(err.Error())
	}

}

// Main Function ( create connection to db, create mux router)

func main(){
// Create MYSQL connection
	db, err = sql.Open ("mysql", "hotdog:hotdog@tcp(127.0.0.1:3306)/hotdog")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

// Create router MUX
	router := mux.NewRouter()
	router.HandleFunc("/hotdogs", getHotdogs).Methods("GET")
	router.HandleFunc("/hotdogs", createHotdog).Methods("POST")
//	router.HandleFunc("/hotdogs/{id}", getHotdog).Methods("GET")
	router.HandleFunc("/hotdogs/{name}", updateHotdog).Methods("PUT")
	router.HandleFunc("/hotdogs/{name}", deleteHotdog).Methods("DELETE")


	fmt.Println("Server started on port :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
