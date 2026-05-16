package main

import (
	"errors"
	"maps"
	"slices"

	"github.com/JordanllHarper/trainsgo/shared"
)

type (
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
		UpsertTrain(t shared.Train) (exists bool, err error)
	}

	trainGetter interface {
		GetTrains() ([]shared.Train, error)
		GetTrainByRef(ref string) (exists bool, t shared.Train, err error)
	}

	inMemTrainStore map[string]shared.Train
)

func (ts inMemTrainStore) GetTrains() ([]shared.Train, error) {
	return slices.Collect(maps.Values(ts)), nil
}

func (ts inMemTrainStore) GetTrainByRef(ref string) (exists bool, t shared.Train, err error) {
	t, exists = ts[ref]
	return exists, t, nil
}

func (ts inMemTrainStore) UpsertTrain(t shared.Train) (exists bool, err error) {
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
