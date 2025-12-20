package main

import (
	"errors"
	"maps"
	"slices"
)

type (
	// A path to take from a start node to an end node
	route struct {
		Id        string           `json:"id"`
		StartNode routeStationNode `json:"start"`
		EndNode   routeStationNode `json:"end"`
	}

	// Node on a graph in a route that is a station a train would stop at
	routeStationNode struct {
		Id           string  `json:"id"`
		RouteId      string  `json:"routeId"`
		StationName  string  `json:"stationId"`
		PreviousNode *string `json:"previous"`
		NextNode     *string `json:"next"`
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
		routeBuilder
	}

	routeBuilder interface {
		// Builds a route from a list of stations
		Build(stations []station) (route, error)
	}
)

type inMemoryRouteStore map[string]route

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
	s[r.Id] = r
	return nil
}

func (s inMemoryRouteStore) Delete(id string) error {
	delete(s, id)
	return nil
}

func Build(stations []station) (route, error) {
	// TODO: Implement
	return route{}, errors.New("TODO: This is complicated :)")
}
