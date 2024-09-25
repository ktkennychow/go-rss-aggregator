package main

import (
	"context"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := strings.TrimPrefix(r.Header.Get("Authorization"), "ApiKey ")
		targetUser, err := cfg.DB.ReadUser(context.Background(), apiKey)
		if err != nil {
			log.Printf("Error reading user from db: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}
		handler(w, r, targetUser)
	}
}