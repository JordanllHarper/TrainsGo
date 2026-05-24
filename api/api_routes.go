package main

import "net/http"

type dependencies struct {
	trainStore
	secretVerifier
}

func setup(d dependencies, srvBase string) http.Handler {
	mux := http.NewServeMux()
	handler := setupTrainRoutes(mux, d, srvBase)
	handler = addLogging(handler)
	return handler
}

func setupTrainRoutes(
	mux *http.ServeMux,
	d dependencies,
	srvBase string,
) http.Handler {

	mux.HandleFunc(getTrains, handleGetTrains(d))
	mux.HandleFunc(getTrainByRef, handleGetTrainByRef(d))
	mux.HandleFunc(getTrainByRefLive, handleGetTrainByRefLive(d))
	mux.HandleFunc(postTrain, handlePostTrain(d, srvBase))
	mux.HandleFunc(patchTrain, addSecretValidation(handlePatchTrain(d, d)))

	return mux
}
