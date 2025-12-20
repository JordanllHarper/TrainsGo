package main

import (
	"fmt"
	"log"
	"net/http"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	buf        []byte
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, []byte{}, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {

	lrw.buf = b
	return lrw.ResponseWriter.Write(b)
}

func addLogging(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loggingWriter := NewLoggingResponseWriter(w)
		method := r.Method
		pathVal := r.URL.Path
		log.Printf("%s request to %s\n", method, pathVal)
		next.ServeHTTP(loggingWriter, r)
		if loggingWriter.statusCode >= 400 {
			log.Printf("%d response from %s - %s\n", loggingWriter.statusCode, pathVal, string(loggingWriter.buf))
			return
		}
		log.Printf("%d response from %s\n", loggingWriter.statusCode, pathVal)
	}
}

func main() {
	// test station store
	stStore := inMemoryStationStore{
		"testStation": Station{
			Name: "testStation",
			Platforms: map[int]bool{
				1: true,
				2: true,
				3: true,
			},
			Neighbors: []string{
				"anotherTest",
				"aThirdTest",
			},
		},
	}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /stations", handleGetStations(stStore))
	mux.HandleFunc("POST /stations", handlePostStations(stStore))
	mux.HandleFunc("PUT /stations/{name}", handlePutStations(stStore))
	mux.HandleFunc("DELETE /stations/{name}", handleDeleteStations(stStore))

	handler := addLogging(mux)

	port := 8080
	log.Println("START:::Listening on port", port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}
