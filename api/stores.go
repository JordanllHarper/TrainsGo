package main

import (
	"maps"
	"slices"
)

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

	routeReader interface {
		ReadAll() ([]route, error)
		ReadById(id string) (exists bool, r route, err error)
	}

	routeWriter interface {
		Insert(r route) error
		Delete(id string) error
	}

	routeStore interface {
		routeReader
		routeWriter
	}

	inMemoryRouteStore   map[string]route
	inMemoryStationStore map[string]station
)

func (s inMemoryStationStore) ReadAll() ([]station, error) {
	return slices.Collect(maps.Values(s)), nil
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

func (s inMemoryRouteStore) ReadAll() ([]route, error) {
	return slices.Collect(maps.Values(s)), nil
}

func (s inMemoryRouteStore) ReadById(id string) (exists bool, r route, err error) {
	v, ok := s[id]
	if !ok {
		return false, route{}, nil
	}

	return true, v, nil
}

func (s inMemoryRouteStore) Insert(r route) error {
	s[r.RouteId] = r
	return nil
}

func (s inMemoryRouteStore) Delete(id string) error {
	delete(s, id)
	return nil
}
