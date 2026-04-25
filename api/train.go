package main

import (
	"errors"
	"maps"
	"slices"
)

type (
	Train struct {
		Ref         string `json:"ref"`
		Description string `json:"description"`
		PosX        int    `json:"posX"`
		PosY        int    `json:"posY"`
	}

	trainStore interface {
		trainGetter
		trainDeleter
		trainUpdater
	}

	trainGetUpdater interface {
		trainUpdater
		trainGetter
	}

	trainDeleter interface {
		DeleteTrain(ref string) error
	}

	trainUpdater interface {
		UpsertTrain(t Train) (exists bool, err error)
	}

	trainGetter interface {
		GetTrains() ([]Train, error)
		GetTrainByRef(ref string) (exists bool, t Train, err error)
	}

	inMemTrainStore map[string]Train
)

func (ts inMemTrainStore) GetTrains() ([]Train, error) {
	return slices.Collect(maps.Values(ts)), nil
}

func (ts inMemTrainStore) GetTrainByRef(ref string) (exists bool, t Train, err error) {
	t, exists = ts[ref]
	return exists, t, nil
}

func (ts inMemTrainStore) UpsertTrain(t Train) (exists bool, err error) {
	if t.Ref == "" {
		return false, errors.New("Invalid empty ref")
	}

	_, exists = ts[t.Ref]
	ts[t.Ref] = t
	return exists, nil
}

func (ts inMemTrainStore) DeleteTrain(ref string) error {
	delete(ts, ref)
	return nil
}
