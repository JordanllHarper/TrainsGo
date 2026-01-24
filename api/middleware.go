package main

import (
	"log"
	"net/http"
)

type middlewareResponseWriter struct {
	http.ResponseWriter
	buf        []byte
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *middlewareResponseWriter {
	return &middlewareResponseWriter{w, []byte{}, http.StatusOK}
}

func (lrw *middlewareResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *middlewareResponseWriter) Write(b []byte) (int, error) {
	lrw.buf = b
	return lrw.ResponseWriter.Write(b)
}

func addLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
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
		})
}
