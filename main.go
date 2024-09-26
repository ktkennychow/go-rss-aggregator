package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/ktkennychow/go-rss-aggregator/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	Queries *database.Queries
	DB *sql.DB
	Ctx context.Context
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func main(){
	godotenv.Load()
	DBURL := os.Getenv("DBURL")
	db, err := sql.Open("postgres", DBURL)
	if err != nil {
		log.Printf("Error connecting to db: %v\n", err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{Queries: dbQueries, DB: db}
	
	serveMux := http.NewServeMux()
	
	serveMux.HandleFunc("GET /v1/healthz", handlerHealthz)
	serveMux.HandleFunc("GET /v1/err", handlerError)

	serveMux.HandleFunc("POST /v1/users", apiCfg.handlerCreateUser)
	serveMux.HandleFunc("GET /v1/users", apiCfg.middlewareAuth(apiCfg.handlerReadUser))

	serveMux.HandleFunc("POST /v1/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	serveMux.HandleFunc("GET /v1/feeds", apiCfg.handlerReadFeeds)
	
	serveMux.HandleFunc("POST /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	serveMux.HandleFunc("GET /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerReadFeedFollows))
	serveMux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))
	
	port := os.Getenv("PORT")
	server := &http.Server{
		Addr: ":" + port,
		Handler: serveMux, 
	}

	log.Printf("Server running on %v\n", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Printf("Error listening on %v: %v\n", port, err)
	}

}