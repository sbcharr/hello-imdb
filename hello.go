package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"fmt"
	//"bytes"
	"compress/gzip"
	"database/sql"
	_ "github.com/lib/pq"
	//"log"

	"github.com/gorilla/mux"
)

type Movie struct {
	Title  string `json:"title"`
	Rating string `json:"rating"`
	Year string `json:"year"`
	IMDBKey string `json:"imdbkey"`
}

const (
	HOST	    = ""
	PORT	    = "5432"
	DB_USER     = ""
	DB_PASSWORD = ""
	DB_NAME     = ""
)



func main() {
	AdminMux := http.NewServeMux()
	router := mux.NewRouter()
	//router.HandleFunc(`/movies`, handleMovies).Methods("GET")
	router.HandleFunc(`/movies/add`, handleMoviesAdd).Methods("PUT")
	AdminMux.Handle("/", router)
	err := http.ListenAndServe(":" + os.Getenv("PORT"), AdminMux)
	if err != nil {
		panic(err)
	}
}


func handleMoviesAdd(res http.ResponseWriter, req *http.Request) {
	//vars := mux.Vars(req)
	//imdbKey := vars["imdbKey"]
	switch req.Method {
	case "PUT":
		var movie []Movie
		unzip, err1 := gzip.NewReader(req.Body)
		if err1 != nil {
			log.Println(err1.Error())
		}

		decoder := json.NewDecoder(unzip)
		error := decoder.Decode(&movie)
		if error != nil {
			log.Println(error.Error())
			http.Error(res, error.Error(), http.StatusInternalServerError)
			return
		}
		dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", HOST, PORT, DB_USER, DB_PASSWORD, DB_NAME)
		db, err := sql.Open("postgres", dbinfo)
		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()

		var sStmt string = "insert into test_app(title, rating, year, imdbkey) values ($1, $2, $3, $4)"
		stmt, err := db.Prepare(sStmt)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		for i := 0; i < len(movie); i++ {
			res, err := stmt.Exec(movie[i].Title, movie[i].Rating, movie[i].Year, movie[i].IMDBKey)
			if err != nil || res == nil {
				log.Fatal(err)
			}
		}

	}
}

/*

func handleMovies(res http.ResponseWriter, req *http.Request) {
	//movies := &Movie{Title:"hello world"}
	res.Header().Set("Content-Type", "application/json")

	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", HOST, PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}
	//pq.checkErr(err)
	defer db.Close()

	rows, err1 := db.Query("SELECT * FROM test_app")
	if err1 != nil {
		log.Fatal(err1)
	}
	for rows.Next() {
		var title string
		var rating string
		var year string
		var imdbKey string
		err = rows.Scan(&title, &rating, &year, &imdbKey)
		if err != nil {
			log.Fatal(err)
		}
		movie := &Movie{Title:title,Rating:rating,Year:year,IMDBKey:imdbKey}
		outputJson, err2 := json.Marshal(movie)
		if err2 != nil {
			log.Println(err2.Error())
			//http.Error(res, error.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(res, string(outputJson))
	}
}

*/


