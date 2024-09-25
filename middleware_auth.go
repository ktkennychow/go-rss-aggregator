package main

import (
	"context"
	"net/http"
	"strings"
)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := strings.TrimPrefix(r.Header.Get("Authorization"), "ApiKey ")
		targetUser, err := cfg.DB.ReadUser(context.Background(), apiKey)
		if err != nil {
			respondWithError(w,http.StatusInternalServerError, err.Error())
		}
		handler(w, r, targetUser)
	}
}