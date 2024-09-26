package main

import (
	"context"
	"net/http"
	"strings"
)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := strings.TrimPrefix(r.Header.Get("Authorization"), "ApiKey ")
		authedUser, err := cfg.DB.ReadUser(context.Background(), apiKey)
		if err != nil {
			respondWithError(w,http.StatusInternalServerError, "User with Api key not found" + err.Error())
			return
		}
		handler(w, r, authedUser)
	}
}