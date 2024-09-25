package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ktkennychow/go-rss-aggregator/internal/database"
)

type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Name string `json:"name"`
}

func (cfg *apiConfig)handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	userToBeCreated := User{}

	dat, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	err = json.Unmarshal(dat, &userToBeCreated)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userToBeCreated.ID = uuid.New()
	userToBeCreated.CreatedAt = time.Now()
	userToBeCreated.UpdatedAt = time.Now()
		
	createdUser, err := cfg.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID: userToBeCreated.ID,
		CreatedAt: userToBeCreated.CreatedAt,
		UpdatedAt: userToBeCreated.UpdatedAt,
		Name: userToBeCreated.Name,
	})
	if err != nil {
		log.Printf("Error creating user in db: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, createdUser)
}

func (cfg *apiConfig)handlerReadUser(w http.ResponseWriter, _ *http.Request, authedUser database.User) {
	respondWithJSON(w, http.StatusOK, authedUser)
}