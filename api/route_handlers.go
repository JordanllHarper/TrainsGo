package main

import (
	"maps"
	"net/http"
	"slices"
)

func handleGetRoutes(s routeReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allRoutes, err := s.ReadAll()
		if err != nil {
			serverError(w, err)
			return
		}

		writeJsonToHttpOk(w, allRoutes)
	}
}

func handleGetRouteById(s routeReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if stringIsEmpty(id) {
			badRequestMsg(w, "Missing id path parameter")
			return
		}

		exists, route, err := s.ReadById(id)
		if err != nil {
			serverError(w, err)
			return
		}

		if !exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		writeJsonToHttpOk(w, route)
	}
}

func handlePostRoute(s stationStore, rs routeStore) http.HandlerFunc {
	type requestDto struct {
		requiredStations map[int]string
	}
	type responseDto struct {
		Id           string  `json:"id"`
		StartStation station `json:"startStation"`
		EndStation   station `json:"endStation"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var dto requestDto
		if err := jsonDecode(r.Body, &dto); err != nil {
			badRequest(w, err)
			return
		}

		if dto.requiredStations == nil {
			badRequestMsg(w, "Missing required stations")
			return
		}

		if len(dto.requiredStations) < 2 {
			badRequestMsg(w, "Required stations list too small. 2 required")
			return
		}

		order := slices.Collect(maps.Keys(dto.requiredStations))
		slices.Sort(order)

		orderedStations := []string{}
		for _, k := range order {
			orderedStations = append(orderedStations, dto.requiredStations[k])
		}

		validNames, stations, err := s.ReadMany(orderedStations)
		if err != nil {
			serverError(w, err)
			return
		}

		if !validNames {
			badRequestMsg(w, "Invalid provided names")
			return
		}
		route, err := buildRoutes(stations)
		if err != nil {
			serverError(w, err)
			return
		}

		if err := rs.Insert(route); err != nil {
			serverError(w, err)
			return
		}

		writeJsonToHttp(w, http.StatusCreated, route)
	}
}

func handleDeleteRoute(s routeWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if stringIsEmpty(id) {
			badRequestMsg(w, "Missing name path parameter")
			return
		}
		if err := s.Delete(id); err != nil {
			serverError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
