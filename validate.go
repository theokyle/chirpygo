package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(params.Body, " ")
	for i, word := range words {
		if _, ok := profaneWords[strings.ToLower(word)]; ok {
			words[i] = "****"
		}

	}
	cleanWords := strings.Join(words, " ")

	respondWithJSON(w, http.StatusOK, returnVals{CleanedBody: cleanWords})
}
