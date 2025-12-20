package main

import (
	"maps"
	"slices"
)

type station struct {
	// unique identifier of the station as well as descriptive name
	Name string `json:"name"`
	// set of platform numbers
	Platforms map[int]bool `json:"platforms"`

	Neighbors []string `json:"neighbors"`
}

type (
	stationReader interface {
		ReadAll() ([]station, error)
		ReadMany(names []string) (validNames bool, sts []station, err error)
		ReadByName(name string) (exists bool, st station, err error)
	}
	stationWriter interface {
		Upsert(st station) (isUpdate bool, err error)
		Delete(name string) error
	}

	stationStore interface {
		stationReader
		stationWriter
	}
)

type inMemoryStationStore map[string]station

func (s inMemoryStationStore) ReadAll() ([]station, error) {
	vals := maps.Values(s)
	return slices.Collect(vals), nil
}

func (s inMemoryStationStore) ReadMany(names []string) (validNames bool, sts []station, err error) {
	stations := []station{}

	for _, name := range names {
		if val, ok := s[name]; ok {
			stations = append(stations, val)
		} else {
			return false, nil, nil
		}
	}

	return true, stations, nil
}
func (s inMemoryStationStore) ReadByName(name string) (exists bool, st station, err error) {
	if val, ok := s[name]; ok {
		return true, val, nil
	}

	return false, station{}, nil
}

func (s inMemoryStationStore) Upsert(st station) (isUpdate bool, err error) {
	_, contains := s[st.Name]
	s[st.Name] = st
	return contains, nil
}

func (s inMemoryStationStore) Delete(name string) error {
	delete(s, name)
	return nil
}
