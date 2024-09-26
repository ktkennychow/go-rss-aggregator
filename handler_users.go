package main

import (
	"context"
	"encoding/json"
	"io"
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
		respondWithError(w,http.StatusInternalServerError, err.Error())
		return
	}
	
	err = json.Unmarshal(dat, &userToBeCreated)
	if err != nil {
		respondWithError(w,http.StatusInternalServerError, err.Error())
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
		respondWithError(w,http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, createdUser)
}

func (cfg *apiConfig)handlerReadUser(w http.ResponseWriter, _ *http.Request, authedUser database.User) {
	respondWithJSON(w, http.StatusOK, authedUser)
}