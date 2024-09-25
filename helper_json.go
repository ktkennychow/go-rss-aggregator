package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with a 5XX err: %v", msg)
	}
	type errorResp struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResp{Error: msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content type", "application/json")
	
	dat, err :=json.Marshal(payload)
	if err != nil {
		respondWithError(w,http.StatusInternalServerError, err.Error())
	}
	
	w.WriteHeader(code)
	w.Write([]byte(dat))
}
