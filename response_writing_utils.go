package main

import "net/http"

func writeJsonToHttp(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return jsonEncode(w, v)
}
func writeJsonToHttpOk(w http.ResponseWriter, v any) error {
	return writeJsonToHttp(w, http.StatusOK, v)
}

func serverError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func badRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}
