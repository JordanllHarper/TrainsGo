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
func writeJsonToHttpOk(w http.ResponseWriter, v any) {
	writeJsonToHttp(w, http.StatusOK, v)
}

func serverError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func badRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func badRequestMsg(w http.ResponseWriter, errMsg string) {
	http.Error(w, errMsg, http.StatusBadRequest)
}
