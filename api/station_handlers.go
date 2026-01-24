package main

import (
	"errors"
	"maps"
	"net/http"
	"slices"
)

func handleGetStations(s stationReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type responseDto struct {
			Name      string   `json:"name"`
			Platforms []int    `json:"platforms"`
			Neighbors []string `json:"neighbors"`
		}

		allStations, err := s.ReadAll()
		dtos := []responseDto{}

		for _, st := range allStations {

			platforms := slices.Sorted(maps.Keys(st.Platforms))
			dtos = append(dtos, responseDto{
				Name:      st.Name,
				Platforms: platforms,
				Neighbors: st.Neighbors,
			})
		}

		if err != nil {
			serverError(w, err)
			return
		}

		writeJsonToHttpOk(w, dtos)
	}
}

func handleGetStationByName(s stationReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type responseDto struct {
			Name      string   `json:"name"`
			Platforms []int    `json:"platforms"`
			Neighbors []string `json:"neighbors"`
		}

		name := r.PathValue("name")
		if stringIsEmpty(name) {
			badRequestMsg(w, "Missing name path parameter")
			return
		}

		exists, st, err := s.ReadByName(name)
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err != nil {
			serverError(w, err)
			return
		}

		platforms := slices.Sorted(maps.Keys(st.Platforms))
		dto := responseDto{
			Name:      st.Name,
			Platforms: platforms,
			Neighbors: st.Neighbors,
		}

		writeJsonToHttpOk(w, dto)
	}
}

func handlePostStations(s stationStore) http.HandlerFunc {
	type requestDto struct {
		Name      *string  `json:"name"`
		Platforms []int    `json:"platforms"`
		Neighbors []string `json:"neighbors"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var dto requestDto
		if err := jsonDecode(r.Body, &dto); err != nil {
			badRequestMsg(w, "Invalid json")
			return
		}
		if stringIsNilOrEmpty(dto.Name) {
			badRequestMsg(w, "Invalid Station Name")
			return
		}

		if dto.Platforms == nil {
			badRequestMsg(w, "Invalid Station Platforms")
			return
		}

		if dto.Neighbors == nil {
			badRequestMsg(w, "Invalid Station Neighbors")
			return
		}

		platforms := map[int]struct{}{}
		for _, v := range dto.Platforms {
			platforms[v] = struct{}{}
		}

		st := station{
			Name:      *dto.Name,
			Platforms: platforms,
			Neighbors: dto.Neighbors,
		}

		exists, _, err := s.ReadByName(*dto.Name)
		if exists {
			badRequestMsg(w, "Already exists")
			return
		}
		if err != nil {
			serverError(w, err)
			return
		}
		if _, err := s.Upsert(st); err != nil {
			serverError(w, err)
			return
		}
		writeJsonToHttp(w, http.StatusCreated, dto)
	}
}

func handlePutStations(s stationWriter) http.HandlerFunc {
	type requestDto struct {
		Platforms []int    `json:"platforms"`
		Neighbors []string `json:"neighbors"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if stringIsEmpty(name) {
			badRequest(w, errors.New("Missing name path parameter"))
			return
		}
		var dto requestDto
		if err := jsonDecode(r.Body, &dto); err != nil {
			badRequest(w, err)
			return
		}

		if dto.Platforms == nil {
			badRequestMsg(w, "Invalid Station Platforms")
		}

		if dto.Neighbors == nil {
			badRequestMsg(w, "Invalid Station Neighbors")
			return
		}

		platforms := map[int]struct{}{}
		for _, v := range dto.Platforms {
			platforms[v] = struct{}{}
		}
		st := station{
			Name:      name,
			Platforms: platforms,
			Neighbors: dto.Neighbors,
		}

		isUpdate, err := s.Upsert(st)
		if err != nil {
			serverError(w, err)
			return
		}
		if isUpdate {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		writeJsonToHttp(w, http.StatusCreated, st)
	}
}

func handleDeleteStations(s stationWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if stringIsEmpty(name) {
			badRequest(w, errors.New("Missing name path parameter"))
			return
		}
		if err := s.Delete(name); err != nil {
			serverError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
