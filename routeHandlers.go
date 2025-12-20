package main

import (
	"errors"
	"log"
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

		if err != nil {
			serverError(w, err)
			return
		}

		if err := writeJsonToHttpOk(w, allRoutes); err != nil {
			log.Println(err)
		}
	}
}

func handleGetRouteById(s routeReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if stringIsEmpty(id) {
			badRequest(w, errors.New("Missing id path parameter"))
			return
		}

		exists, route, err := s.ReadById(id)
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err != nil {
			serverError(w, err)
			return
		}

		if err := writeJsonToHttpOk(w, route); err != nil {
			log.Println(err)
		}
	}
}

func handlePostRoute(s stationStore, rs routeStore, rb routeBuilder) http.HandlerFunc {
	type requestDto struct {
		requiredStations map[int]string
	}
	type responseDto struct {
		Id           string  `json:"id"`
		StartStation station `json:"startStation"`
		EndStation   station `json:"endStation"`
	}
	validateInput := func(dto requestDto) error {
		if dto.requiredStations == nil {
			return errors.New("Missing required stations")
		}

		if len(dto.requiredStations) < 2 {
			return errors.New("Required stations list too small. 2 required")
		}

		return nil
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var reqDto requestDto
		if err := jsonDecode(r.Body, &reqDto); err != nil {
			badRequest(w, err)
			return
		}
		err := validateInput(reqDto)
		if err != nil {
			badRequest(w, err)
			return
		}

		order := slices.Collect(maps.Keys(reqDto.requiredStations))
		slices.Sort(order)
		orderedStations := []string{}
		for _, k := range order {
			orderedStations = append(orderedStations, reqDto.requiredStations[k])
		}

		validNames, stations, err := s.ReadMany(orderedStations)
		if err != nil {
			serverError(w, err)
			return
		}

		if !validNames {
			badRequest(w, errors.New("Invalid provided names"))
			return
		}
		route, err := rb.Build(stations)
		if err != nil {
			serverError(w, err)
			return
		}

		if err := rs.Insert(route); err != nil {
			serverError(w, err)
			return
		}

		if err := writeJsonToHttp(w, http.StatusCreated, route); err != nil {
			log.Println(err)
		}
	}
}

func handleDeleteRoute(s routeWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if stringIsEmpty(id) {
			badRequest(w, errors.New("Missing name path parameter"))
			return
		}
		if err := s.Delete(id); err != nil {
			serverError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
