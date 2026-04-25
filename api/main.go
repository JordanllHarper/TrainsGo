package main

import (
	"log"
	"net/http"
)

func main() {
	addr := "127.0.0.1:8080"
	stores := setupMock()
	handler := setup(stores, addr)
	log.Println("START:::Listening at", addr)
	log.Fatalln(http.ListenAndServe(addr, handler))
}
func setupMock() dependencies {
	return dependencies{
		ts: inMemTrainStore{},
	}
}
