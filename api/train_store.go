package main

import (
	"errors"
	"maps"
	"slices"
	"strings"

	"github.com/JordanllHarper/trainsgo/shared"
)

var (
	errorEmptyRef      = errors.New("invalid empty ref")
	errorAlreadyExists = errors.New("train already exists")
)

type (
	trainStore interface {
		trainGetter
		trainDeleter
		trainCreater
		trainUpdater
	}

	trainGetUpdater interface {
		trainUpdater
		trainGetter
	}

	trainDeleter interface {
		DeleteTrain(ref string) error
	}

	trainCreater interface {
		CreateTrain(t shared.Train) error
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

func (ts inMemTrainStore) CreateTrain(t shared.Train) error {
	if strings.TrimSpace(t.Ref) == "" {
		return errorEmptyRef
	}
	if _, exists := ts[t.Ref]; exists {
		return errorAlreadyExists
	}
	ts[t.Ref] = t
	return nil
}

func (ts inMemTrainStore) UpsertTrain(t shared.Train) (exists bool, err error) {
	if strings.TrimSpace(t.Ref) == "" {
		return false, errorEmptyRef
	}

	_, exists = ts[t.Ref]
	ts[t.Ref] = t
	return exists, nil
}

func (ts inMemTrainStore) DeleteTrain(ref string) error {
	delete(ts, ref)
	return nil
}
