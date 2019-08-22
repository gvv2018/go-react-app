package main

import (
	"log"
	"net/http"
	"database/sql"
	"encoding/json"
	"fmt"
//	"io/ioutil"
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


func getHotdogs( w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var hotdogs []Hotdog

	result, err := db.Query("SELECT name, price FROM type_hotdogs")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

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

func main(){
// Create MYSQL connection
	db, err = sql.Open ("mysql", "hotdog:hotdog@tcp(127.0.0.1:3306)/hotdog")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

// Create router MUX
	router := mux.NewRouter()
//	hotdogs = append(hotdogs, Hotdog{Name: "Super", Price: "24"})
//	hotdogs = append(hotdogs, Hotdog{Name: "Super 2", Price: "15"})
	router.HandleFunc("/hotdogs", getHotdogs).Methods("GET")


	fmt.Println("Server started on port :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
