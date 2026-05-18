package main

import (
	"log"
	"net/http"

	"github.com/JordanllHarper/trainsgo/shared"
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
		trainStore: inMemTrainStore{
			"test": shared.Train{
				Ref: "test",
			},
		},
		secretVerifier: inMemSecretVerifier{
			"test": "test_key",
		},
	}
}
