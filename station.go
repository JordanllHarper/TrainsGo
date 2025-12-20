package main

import (
	"maps"
	"slices"
)

type Station struct {
	// unique identifier of the station as well as descriptive name
	Name string `json:"name"`
	// set of platform numbers
	Platforms map[int]bool `json:"platforms"`

	Neighbors []string `json:"neighbors"`
}

type (
	stationReader interface {
		ReadAll() ([]Station, error)
		ReadByName(name string) (exists bool, st Station, err error)
	}
	stationWriter interface {
		Upsert(st Station) (isUpdate bool, err error)
		Delete(name string) error
	}

	stationStore interface {
		stationReader
		stationWriter
	}
)

type inMemoryStationStore map[string]Station

func (s inMemoryStationStore) ReadAll() ([]Station, error) {
	vals := maps.Values(s)
	return slices.Collect(vals), nil
}
func (s inMemoryStationStore) ReadByName(name string) (exists bool, st Station, err error) {
	if val, ok := s[name]; ok {
		return true, val, nil
	}

	return false, Station{}, nil
}

func (s inMemoryStationStore) Upsert(st Station) (isUpdate bool, err error) {
	_, contains := s[st.Name]
	s[st.Name] = st
	return contains, nil
}

func (s inMemoryStationStore) Delete(name string) error {
	delete(s, name)
	return nil
}
