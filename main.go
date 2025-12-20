package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	stores := setupMock()
	handler := setup(stores)
	ip := "127.0.0.1"
	port := "8080"
	log.Println("START:::Listening on port", port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), handler))
}
