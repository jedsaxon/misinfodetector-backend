package main

import (
	"log"
	"misinfodetector-backend/config"
	"misinfodetector-backend/dbservice"
	"misinfodetector-backend/handler"
	"misinfodetector-backend/handler/middleware"
	"misinfodetector-backend/mqservice"
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
	dbs, sqliteClose, err := dbservice.NewDbService(cfg.SqliteDsn)
	if err != nil {
		log.Fatalf("error opening sqlite database: %v", err)
	}
	defer func() {
		err := sqliteClose()
		if err != nil {
			log.Panicf("error closing database connection: %v", err)
		}
	}()

	// Connect to RabbitMQ
	log.Printf("connecting and initialising rabbitmq")
	mqs, mqClose, err := mqservice.NewMqService(&mqservice.MqServiceConfig{
		Uri:             cfg.RabbitMqUri,
		InputQueueName:  cfg.RabbitMqInputQueueName,
		OutputQueueName: cfg.RabbitMqOutputQueueName,
	})
	if err != nil {
		log.Fatalf("error opening rabbitmq service: %v", err)
	}
	defer func() {
		err := mqClose()
		if err != nil {
			log.Panicf("error closing rabbitmq connection: %v", err)
		}
	}()

	// Create Posts Controller
	c := handler.NewPostsController(dbs, mqs)

	// Configure RabbitMQ Queues
	go mqs.SubscribeToMisinfoOutput(c.HandleNewMisinfoReport)

	// Configure Router
	log.Printf("configuring router")
	r := mux.NewRouter()

	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.ContentTypeJsonMiddleware)

	r.HandleFunc("/api/posts", c.GetPosts).Methods(http.MethodGet)
	r.HandleFunc("/api/posts", c.PutPosts).Methods(http.MethodPut)
	r.HandleFunc("/api/posts", c.UploadPost).Methods(http.MethodPost)
	r.HandleFunc("/api/posts/all", c.GetAllPosts).Methods(http.MethodGet)
	r.HandleFunc("/api/posts/{id}", c.GetSpecificPost).Methods(http.MethodGet)
	r.HandleFunc("/api/posts/random", c.PutRandomPosts).Methods(http.MethodPost)

	r.HandleFunc("/api/data/tnse-embeddings", c.GetTnseEmbeddings).Methods(http.MethodGet)
	r.HandleFunc("/api/data/tnse-embeddings", c.PutTnseEmbeddings).Methods(http.MethodPut)
	r.HandleFunc("/api/data/topic-activities", c.GetTopicActivities).Methods(http.MethodGet)
	r.HandleFunc("/api/data/topic-activities", c.PutTopicActivities).Methods(http.MethodPut)

	handler := cors.AllowAll().Handler(r)

	// Serve
	log.Printf("listening on %s", cfg.ListenAddres)
	err = http.ListenAndServe(cfg.ListenAddres, handler)
	if err != nil {
		log.Fatalf("error while listening for requests: %v", err)
	}
}
