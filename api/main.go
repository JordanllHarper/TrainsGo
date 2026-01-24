package main

import (
	"log"
	"net/http"
)

func main() {
	stores := setupMock()
	handler := setup(stores)
	addr := "127.0.0.1:8080"
	log.Println("START:::Listening at", addr)
	log.Fatalln(http.ListenAndServe(addr, handler))
}
