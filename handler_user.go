package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
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

func (apiConfig apiConfig)handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	newUser := User{}

	dat, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	err = json.Unmarshal(dat, &newUser)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newUser.ID = uuid.New()
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()
		
	createdUser, err := apiConfig.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID: newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Name: newUser.Name,
	})
	if err != nil {
		log.Printf("Error creating user in db: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, createdUser)
}

func (apiConfig apiConfig)handlerReadUser(w http.ResponseWriter, r *http.Request) {
	apiKey := strings.TrimPrefix(r.Header.Get("Authorization"), "ApiKey ")

	targetUser, err := apiConfig.DB.ReadUser(context.Background(), apiKey)
	if err != nil {
		log.Printf("Error reading user from db: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, targetUser)
}