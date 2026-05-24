package main

import (
	"log"
	"net/http"

	"github.com/JordanllHarper/trainsgo/shared"
)

type dependencies struct {
	trainStore
	stationStore
	secretVerifier
}

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
			listeners: map[string][]listener{},
			trains: map[string]shared.Train{
				"test": {
					Ref: "test",
				},
			},
		},
		stationStore: inMemStationStore{},
		secretVerifier: inMemSecretVerifier{
			"test": "test_key",
		},
	}
}

func setup(d dependencies, srvBase string) http.Handler {
	mux := http.NewServeMux()
	handler := setupTrainRoutes(mux, d, srvBase)
	handler = setupStationRoutes(mux, d, srvBase)
	handler = addLogging(handler)
	return handler
}
