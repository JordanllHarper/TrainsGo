package main

import (
	"log"
	"net/http"
)

func writeJsonToHttp(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := jsonEncode(w, v); err != nil {
		log.Println(err)
	}
}

func writeJsonToHttpCreated(w http.ResponseWriter, location string, v any) {
	w.Header().Add("Location", location)
	writeJsonToHttp(w, http.StatusCreated, v)
}
func writeJsonToHttpOk(w http.ResponseWriter, v any) {
	writeJsonToHttp(w, http.StatusOK, v)
}

func internalServerError(w http.ResponseWriter, err error) {
	log.Println("Error:", err)
	http.Error(w, "Something went wrong", http.StatusInternalServerError)
}

func badRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}
