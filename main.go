package main

import (
	"database/sql"
	"log"
	"misinfodetector-backend/dbservice"
	"misinfodetector-backend/handler"
	"misinfodetector-backend/handler/middleware"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"

	"github.com/gorilla/mux"
)

func main() {
	log.Printf("connecting and initialising sqlite")
	db, err := sql.Open("sqlite3", "file:/var/lib/backend/app.db?cache=shared")
	if err != nil {
		log.Fatalf("error opening sqlite database: %v", err)
	}
	defer db.Close()
	dbs := dbservice.NewDbService(db)

	c := handler.NewPostsController(dbs)
	r := mux.NewRouter()

	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.ContentTypeJsonMiddleware)
	r.HandleFunc("/api/posts", c.GetPosts).Methods(http.MethodGet)
	r.HandleFunc("/api/posts", c.PutPost).Methods(http.MethodPost)

	handler := cors.AllowAll().Handler(r)

	listen := "0.0.0.0:5000"
	log.Printf("listening on %s", listen)
	err = http.ListenAndServe(listen, handler)
	if err != nil {
		log.Fatalf("error while listening for requests: %v", err)
	}
}
