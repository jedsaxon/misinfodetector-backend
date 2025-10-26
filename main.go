package main

import (
	"database/sql"
	"log"
	"misinfodetector-backend/config"
	"misinfodetector-backend/dbservice"
	"misinfodetector-backend/handler"
	"misinfodetector-backend/handler/middleware"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"

	"github.com/gorilla/mux"
)

func main() {
	// Load Environment Variables
	cfg := config.NewDefaultConfig()
	cfg.PopulateFromEnv()
	cfg.PopulateFromArgs()

	// Connect to Sqlite
	log.Printf("connecting and initialising sqlite")
	db, err := sql.Open("sqlite3", cfg.SqliteDsn)
	if err != nil {
		log.Fatalf("error opening sqlite database: %v", err)
	}
	defer db.Close()
	dbs := dbservice.NewDbService(db)

	// Configure Router
	c := handler.NewPostsController(dbs)
	r := mux.NewRouter()

	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.ContentTypeJsonMiddleware)
	r.HandleFunc("/api/posts", c.GetPosts).Methods(http.MethodGet)
	r.HandleFunc("/api/posts", c.PutPost).Methods(http.MethodPost)
	r.HandleFunc("/api/posts/random", c.PutRandomPost).Methods(http.MethodPost)

	handler := cors.AllowAll().Handler(r)

	// Serve 
	log.Printf("listening on %s", cfg.ListenAddres)
	err = http.ListenAndServe(cfg.ListenAddres, handler)
	if err != nil {
		log.Fatalf("error while listening for requests: %v", err)
	}
}
