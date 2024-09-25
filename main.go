package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/ktkennychow/go-rss-aggregator/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func main(){
	godotenv.Load()
	DBURL := os.Getenv("DBURL")
	DB, err := sql.Open("postgres", DBURL)
	if err != nil {
		log.Printf("Error connecting to db: %v\n", err)
	}

	dbQueries := database.New(DB)

	apiCfg := apiConfig{DB: dbQueries}
	
	serveMux := http.NewServeMux()
	
	serveMux.HandleFunc("GET /v1/healthz", handlerHealthz)
	serveMux.HandleFunc("GET /v1/err", handlerError)
	serveMux.HandleFunc("POST /v1/users", apiCfg.handlerCreateUser)
	serveMux.HandleFunc("GET /v1/users", apiCfg.middlewareAuth(apiCfg.handlerReadUser))
	serveMux.HandleFunc("POST /v1/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	
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