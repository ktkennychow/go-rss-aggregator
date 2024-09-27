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

type Feed struct {
	ID uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Name string `json:"name"`
	Url string `json:"url"`
	UserID uuid.UUID `json:"user_id"`
	LastFetchedAt *time.Time `json:"last_fetched_at"`
}

func (cfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, authedUser database.User) {
	feedToBeCreated := Feed{}

	dat, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w,http.StatusInternalServerError, "Error reading request body:" + err.Error())
		return
	}
	
	err = json.Unmarshal(dat, &feedToBeCreated)
	if err != nil {
		respondWithError(w,http.StatusInternalServerError, "Error unmarshaling request body:" + err.Error())
		return
	}

	feedToBeCreated.ID = uuid.New()

	tx, err := cfg.DB.Begin()
	if err != nil {
		respondWithError(w,http.StatusInternalServerError, "error creating feed transaction:" + err.Error())
		return
	}
	
	qtx := cfg.Queries.WithTx(tx)
		
	createdFeed, err := qtx.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: feedToBeCreated.ID,
		Name: feedToBeCreated.Name,
		Url: feedToBeCreated.Url,
		UserID: authedUser.ID,
	})
	if err != nil {
		tx.Rollback()
		respondWithError(w,http.StatusInternalServerError, "failed to create a feed in db: " + err.Error())
		return
	}

	newFeedFollowID := uuid.New()

	createdFeedFollow, err := qtx.CreateUserFeed(context.Background(), database.CreateUserFeedParams{
		ID: newFeedFollowID,
		FeedID: feedToBeCreated.ID,
		UserID: authedUser.ID,
	})
	if err != nil {
		tx.Rollback()
		respondWithError(w, http.StatusInternalServerError, "failed to create a feed follow in db: " + err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not commit feed transaction:" + err.Error())
		return
	}

	type ResBody struct {
		Feed Feed `json:"feed"`
		FeedFollow FeedFollow `json:"feed_follow"`
	}

	var normalizedLastFetchAt *time.Time
	if (createdFeed.LastFetchedAt.Valid) {
		normalizedLastFetchAt = &createdFeed.LastFetchedAt.Time
	}

	resBody := ResBody{
		Feed: Feed{
			ID: createdFeed.ID,
			CreatedAt: createdFeed.CreatedAt,
			UpdatedAt: createdFeed.UpdatedAt,
			Name: createdFeed.Name,
			Url: createdFeed.Url,
			UserID: authedUser.ID,
			LastFetchedAt: normalizedLastFetchAt,
		},
		FeedFollow: FeedFollow{
			ID: createdFeedFollow.ID,
			CreatedAt: createdFeedFollow.CreatedAt,
			UpdatedAt: createdFeedFollow.UpdatedAt,
			FeedID: createdFeedFollow.FeedID,
			UserID: authedUser.ID,
		},
	}

	respondWithJSON(w, http.StatusOK, resBody)
}

func (cfg *apiConfig) handlerReadFeeds(w http.ResponseWriter, _ *http.Request) {
	feeds, err := cfg.Queries.ReadFeeds(context.Background())
	if err != nil {
		respondWithError(w,http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, feeds)
}