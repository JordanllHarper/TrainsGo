package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/JordanllHarper/trainsgo/shared"
)

func setupStationRoutes(
	mux *http.ServeMux,
	d dependencies,
	srvBase string,
) http.Handler {

	mux.HandleFunc(getStations, handleGetStations(d))
	mux.HandleFunc(getStationById, handleGetStationById(d))
	mux.HandleFunc(postStation, handlePostStation(d, srvBase))
	mux.HandleFunc(patchStation, handlePatchStationById(d))
	mux.HandleFunc(deleteStation, handleDeleteStationById(d))

	return mux
}

const getStations = "GET /stations"

func handleGetStations(sg stationGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stations, err := sg.GetStations()
		if err != nil {
			internalServerError(w, err)
			return
		}
		writeJsonToHttpOk(w, stations)
	}
}

const getStationById = "GET /stations/{id}"

func handleGetStationById(sg stationGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ref := r.PathValue("id")
		exists, station, err := sg.GetStationById(ref)
		if err != nil {
			internalServerError(w, err)
			return
		}

		if !exists {
			http.NotFound(w, r)
			return
		}

		writeJsonToHttpOk(w, station)
	}
}

const postStation = "POST /stations"

func handlePostStation(tu stationCreator, srvBase string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var createReq shared.CreateStationRequest
		if err := jsonDecode(r.Body, &createReq); err != nil {
			badRequest(w, err)
			return
		}

		s, err := tu.CreateStation(createReq.Name, createReq.PosX, createReq.PosY)
		switch {
		case errors.Is(err, errorEmptyStationId):
			badRequest(w, err)
			return
		case errors.Is(err, errorAlreadyExists):
			http.Error(w, "Station with id already exists", http.StatusConflict)
			return
		}

		writeJsonToHttpCreated(
			w,
			fmt.Sprintf("%v/stations/%v", srvBase, s.Id),
			s,
		)
	}
}

const patchStation = "PATCH /stations/{id}"

func handlePatchStationById(su stationUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		var req shared.PatchStationRequest
		if err := jsonDecode(r.Body, &req); err != nil {
			badRequest(w, err)
			return
		}

		t, err := su.UpdateStation(
			id,
			req.Name,
		)

		switch {
		case errors.Is(err, errorEmptyStationId):
			badRequest(w, err)
			return
		case errors.Is(err, errorNotFound):
			http.NotFound(w, r)
			return
		case err != nil:
			internalServerError(w, err)
			return
		}

		writeJsonToHttpOk(w, t)
	}
}

const deleteStation = "DELETE /stations/{id}"

func handleDeleteStationById(su stationDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		if err := su.DeleteStation(id); err != nil {
			internalServerError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
