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

type FeedFollow struct {
	ID uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	UserID uuid.UUID `json:"user_id"`
	FeedID uuid.UUID `json:"feed_id"`
}

func (cfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, authedUser database.User) {
	feedFollowToBeCreated := FeedFollow{}
	dat, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.
		StatusInternalServerError, err.Error())
		return
	}
	
	err = json.Unmarshal(dat, &feedFollowToBeCreated)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	feedFollowToBeCreated.ID = uuid.New()
	feedFollowToBeCreated.CreatedAt = time.Now()
	feedFollowToBeCreated.UpdatedAt = time.Now()
		
	createdFeedFollow, err := cfg.Queries.CreateUserFeed(context.Background(), database.CreateUserFeedParams{
		ID: feedFollowToBeCreated.ID,
		CreatedAt: feedFollowToBeCreated.CreatedAt,
		UpdatedAt: feedFollowToBeCreated.UpdatedAt,
		FeedID: feedFollowToBeCreated.FeedID,
		UserID: authedUser.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, createdFeedFollow)
}

func (cfg *apiConfig) handlerReadFeedFollows(w http.ResponseWriter, _ *http.Request, authedUser database.User) {
	feedfollows, err := cfg.Queries.ReadUserFeeds(context.Background(), authedUser.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, feedfollows)
}

func (cfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, authedUser database.User) {
	feedFollowIDString := r.PathValue("feedFollowID")
	
	feedFollowID, err := uuid.Parse(feedFollowIDString)
	if err != nil {
		respondWithError(w,http.StatusInternalServerError, err.Error())
		return
	}
	
	deletedFeedFollow, err := cfg.Queries.DeleteUserFeed(context.Background(), database.DeleteUserFeedParams{
		ID: feedFollowID,
		UserID: authedUser.ID,
	})
	if err != nil {
		respondWithError(w,http.StatusInternalServerError, "Failed to Delete feed follow in db" + err.Error())
		return 
	}

	respondWithJSON(w, http.StatusOK, deletedFeedFollow)
}