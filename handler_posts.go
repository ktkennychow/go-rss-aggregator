package main

import (
	"net/http"
	"strconv"

	"github.com/ktkennychow/go-rss-aggregator/internal/database"
)

func (cfg *apiConfig) handlerReadPostsByUser(w http.ResponseWriter, r *http.Request, authedUser database.User) {
	limitParamStr := r.URL.Query().Get("limit")
	limitParamInt, err := strconv.Atoi(limitParamStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	limit := 5
	if limitParamInt > 0 {
		limit = limitParamInt
	}

	posts, err := cfg.Queries.ReadPostsByUser(cfg.Ctx,database.ReadPostsByUserParams{
		UserID: authedUser.ID,
		Limit: int32(limit),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusOK, posts)
}