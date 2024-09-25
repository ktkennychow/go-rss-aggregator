package main

import "net/http"

func handlerHealthz(w http.ResponseWriter, _ *http.Request) {
	respBody := map[string]string{
		"status": "ok",
	}
	respondWithJSON(w, 200, respBody)
}

func handlerError(w http.ResponseWriter, _ *http.Request) {
	respondWithError(w, 500, "Internal Server Error")
}