package main

import (
	"errors"
	"maps"
	"slices"
)

type (
	// A path to take from a start node to an end node
	route struct {
		RouteId string          `json:"routeId"`
		Route   map[int]station `json:"route"`
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

	routeBuilder interface {
		// Builds a route from a list of stations
		Build(stations []station) (route, error)
	}

	inMemoryRouteStore map[string]route

	inMemoryRouteBuilder map[string]station
)

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

func (b inMemoryRouteBuilder) Build(stations []station) (route, error) {
	// TODO: Implement
	return route{}, errors.New("TODO: This is complicated :)")
}
