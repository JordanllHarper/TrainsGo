package main

import "net/http"

type stores struct {
	ss stationStore
	rs routeStore
}

func setup(stores stores) http.Handler {
	mux := http.NewServeMux()
	handler := setupStationRoutes(mux, stores.ss)
	handler = setupRouteRoutes(mux, stores.rs, stores.ss)
	handler = addLogging(handler)
	return handler
}

func setupStationRoutes(mux *http.ServeMux, ss stationStore) http.Handler {
	mux.HandleFunc("GET /stations", handleGetStations(ss))
	mux.HandleFunc("GET /stations/{name}", handleGetStationByName(ss))
	mux.HandleFunc("POST /stations", handlePostStations(ss))
	mux.HandleFunc("PUT /stations/{name}", handlePutStations(ss))
	mux.HandleFunc("DELETE /stations/{name}", handleDeleteStations(ss))
	return mux
}

func setupRouteRoutes(mux *http.ServeMux, rs routeStore, ss stationStore) http.Handler {
	mux.HandleFunc("GET /routes", handleGetRoutes(rs))
	mux.HandleFunc("POST /routes", handlePostRoute(ss, rs))
	mux.HandleFunc("GET /routes/{id}", handleGetRouteById(rs))
	mux.HandleFunc("DELETE /routes/{id}", handleDeleteRoute(rs))
	return mux
}
