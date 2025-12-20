package main

import (
	"errors"
	"log"
	"maps"
	"net/http"
	"slices"
)

func handleGetStations(s stationReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		type responseDto struct {
			Name      string   `json:"name"`
			Platforms []int    `json:"platforms"`
			Neighbors []string `json:"neighbors"`
		}

		mapToDto := func(st Station) responseDto {
			platforms := slices.Sorted(maps.Keys(st.Platforms))
			return responseDto{
				Name:      st.Name,
				Platforms: platforms,
				Neighbors: st.Neighbors,
			}
		}
		w.Header().Add("content-type", "application/json")
		if !query.Has("name") {
			allStations, err := s.ReadAll()
			dtos := []responseDto{}

			for _, st := range allStations {
				dtos = append(dtos, mapToDto(st))
			}

			if err != nil {
				serverError(w, err)
				return
			}

			if err := jsonEncode(w, dtos); err != nil {
				log.Fatalln(err)
			}
			return
		}

		name := query.Get("name")

		exists, st, err := s.ReadByName(name)
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err != nil {
			serverError(w, err)
			return
		}

		dto := mapToDto(st)

		if err := jsonEncode(w, dto); err != nil {
			log.Println(err)
		}
	}
}

func handlePostStations(s stationStore) http.HandlerFunc {
	type requestDto struct {
		Name      *string  `json:"name"`
		Platforms []int    `json:"platforms"`
		Neighbors []string `json:"neighbors"`
	}
	validate := func(dto requestDto) (Station, error) {
		if stringIsNilOrEmpty(dto.Name) {
			return Station{}, errors.New("Invalid Station Name")
		}

		if dto.Platforms == nil {
			return Station{}, errors.New("Invalid Station Platforms")
		}

		if dto.Neighbors == nil {
			return Station{}, errors.New("Invalid Station Neighbors")
		}

		platforms := map[int]bool{}
		for _, v := range dto.Platforms {
			platforms[v] = true
		}

		return Station{
			Name:      *dto.Name,
			Platforms: platforms,
			Neighbors: dto.Neighbors,
		}, nil

	}
	return func(w http.ResponseWriter, r *http.Request) {
		var reqDto requestDto
		if err := jsonDecode(r.Body, &reqDto); err != nil {
			badRequest(w, err)
			return
		}
		st, err := validate(reqDto)
		if err != nil {
			badRequest(w, err)
			return
		}
		exists, _, err := s.ReadByName(st.Name)
		if exists {
			badRequest(w, errors.New("Already exists"))
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
		w.WriteHeader(http.StatusCreated)
		if err := jsonEncode(w, reqDto); err != nil {
			log.Println(err)
		}
	}
}

func handlePutStations(s stationWriter) http.HandlerFunc {
	type requestDto struct {
		Platforms []int    `json:"platforms"`
		Neighbors []string `json:"neighbors"`
	}
	validate := func(name string, dto requestDto) (Station, error) {
		if dto.Platforms == nil {
			return Station{}, errors.New("Invalid Station Platforms")
		}

		if dto.Neighbors == nil {
			return Station{}, errors.New("Invalid Station Neighbors")
		}

		platforms := map[int]bool{}
		for _, v := range dto.Platforms {
			platforms[v] = true
		}
		return Station{
			Name:      name,
			Platforms: platforms,
			Neighbors: dto.Neighbors,
		}, nil

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

		st, err := validate(name, dto)
		if err != nil {
			badRequest(w, err)
			return
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
		w.WriteHeader(http.StatusCreated)
		if err := jsonEncode(w, st); err != nil {
			log.Println(err)
		}

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

func serverError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func badRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}
