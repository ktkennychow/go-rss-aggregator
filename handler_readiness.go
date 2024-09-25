package main

import "net/http"

func handlerHealthz(w http.ResponseWriter, _ *http.Request) {
	respBody := map[string]string{
		"status": "ok",
	}
	respondWithJSON(w, http.StatusOK, respBody)
}

func handlerError(w http.ResponseWriter, _ *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}