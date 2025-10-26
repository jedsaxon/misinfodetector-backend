package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"misinfodetector-backend/dbservice"
	"misinfodetector-backend/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db, err := sql.Open("sqlite3", "file:app.db?cache=shared&mode=memory")
	if err != nil {
		log.Fatalf("error opening sqlite database: %v", err)
	}
	defer db.Close()
	dbs := dbservice.NewDbService(db)

	log.Printf("initialising sqlite")
	initDb(db)

	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	r.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) { handler.GetPosts(w, r,dbs) }).Methods(http.MethodGet)
	r.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) { handler.PutPost(w, r,dbs) }).Methods(http.MethodPut)

	listen := "0.0.0.0:3000"
	log.Printf("listening on %s", listen)
	err = http.ListenAndServe(listen, r)
	if err != nil {
		log.Fatalf("error while listening for requests: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("request recieved - %s", r.RequestURI)
        next.ServeHTTP(w, r)
    })
}

func initDb(db *sql.DB) error {
	_, err := db.Exec("create table if not exists posts(id varchar(36), message text, username text, date_submitted text, is_misinformation int);")
	if err != nil {
		return err
	}
	return nil
}
