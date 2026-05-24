package main

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/JordanllHarper/trainsgo/shared"
	"github.com/google/uuid"
)

var (
	errorEmptyStationId = errors.New("Invalid empty station id")
)

func withStationContext(err error) error { return fmt.Errorf("station %w", err) }

type (
	stationStore interface {
		stationCreator
		stationGetter
		stationUpdater
		stationDeleter
	}

	stationGetter interface {
		GetStations() ([]shared.Station, error)
		GetStationById(id string) (exists bool, s shared.Station, err error)
	}

	stationCreator interface {
		CreateStation(name string, posX, posY int) (shared.Station, error)
	}

	stationUpdater interface {
		UpdateStation(id string, name string) (shared.Station, error)
	}

	stationDeleter interface {
		DeleteStation(id string) error
	}
	inMemStationStore map[string]shared.Station
)

func (ss inMemStationStore) GetStations() ([]shared.Station, error) {
	return slices.Collect(maps.Values(ss)), nil
}

func (ss inMemStationStore) GetStationById(id string) (exists bool, s shared.Station, err error) {
	s, ok := ss[id]
	return ok, s, nil
}

func (ss inMemStationStore) CreateStation(name string, posX int, posY int) (shared.Station, error) {
	id := uuid.NewString()
	s := shared.Station{
		Id:   id,
		Name: name,
		PosX: posX,
		PosY: posY,
	}
	ss[id] = s
	return s, nil
}

func (ss inMemStationStore) UpdateStation(id string, name string) (shared.Station, error) {
	if strings.TrimSpace(id) == "" {
		return shared.Station{}, errorEmptyStationId
	}
	s, ok := ss[id]
	if !ok {
		return shared.Station{}, withStationContext(errorNotFound)
	}
	s.Name = name
	ss[id] = s
	return s, nil
}

func (ss inMemStationStore) DeleteStation(id string) error {
	delete(ss, id)
	return nil
}
