package main

import "net/http"

type dependencies struct {
	ts trainStore
}

func setup(d dependencies, srvBase string) http.Handler {
	mux := http.NewServeMux()
	handler := setupTrainRoutes(mux, d.ts, srvBase)
	handler = addLogging(handler)
	return handler
}

func setupTrainRoutes(
	mux *http.ServeMux,
	ts trainStore,
	srvBase string,
) http.Handler {

	mux.HandleFunc(getTrains, HandleGetTrains(ts))
	mux.HandleFunc(getTrainByRef, HandleGetTrainByRef(ts))
	mux.HandleFunc(postTrain, HandlePostTrain(ts, srvBase))
	mux.HandleFunc(patchTrain, HandlePatchTrain(ts))

	return mux
}
