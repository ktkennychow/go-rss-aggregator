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

type Feed struct {
	ID uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Name string `json:"name"`
	Url string `json:"url"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, authedUser database.User) {
	feedToBeCreated := Feed{}

	dat, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	err = json.Unmarshal(dat, &feedToBeCreated)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	feedToBeCreated.ID = uuid.New()
	feedToBeCreated.CreatedAt = time.Now()
	feedToBeCreated.UpdatedAt = time.Now()
		
	createdFeed, err := cfg.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: feedToBeCreated.ID,
		CreatedAt: feedToBeCreated.CreatedAt,
		UpdatedAt: feedToBeCreated.UpdatedAt,
		Name: feedToBeCreated.Name,
		Url: feedToBeCreated.Url,
		UserID: authedUser.ID,
	})
	if err != nil {
		respondWithError(w,http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, createdFeed)
}

func (cfg *apiConfig) handlerReadFeeds(w http.ResponseWriter, _ *http.Request) {
	feeds, err := cfg.DB.ReadFeeds(context.Background())
	if err != nil {
		respondWithError(w,http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, feeds)
}