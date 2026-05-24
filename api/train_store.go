package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"maps"
	"slices"
	"strings"

	"github.com/JordanllHarper/trainsgo/shared"
)

var (
	errorEmptyTrainRef = errors.New("invalid empty ref")
)

func withTrainContext(err error) error { return fmt.Errorf("train %w", err) }

type (
	trainStore interface {
		trainGetter
		trainCreator
		trainUpdater
		trainUpdateSender
	}

	trainCreator interface {
		CreateTrain(t shared.Train) error
	}

	trainUpdateSender interface {
		RegisterListener(ref string, ctx context.Context) (chan shared.Train, error)
	}

	trainUpdater interface {
		UpdateTrain(ref string, desc *string, posX, posY *int) (shared.Train, error)
	}

	trainGetter interface {
		GetTrains() ([]shared.Train, error)
		GetTrainByRef(ref string) (exists bool, t shared.Train, err error)
	}

	inMemTrainStore struct {
		trains    map[string]shared.Train
		listeners map[string][]listener
	}

	listener struct {
		sendCh chan shared.Train
		ctx    context.Context
	}
)

func (ts inMemTrainStore) GetTrains() ([]shared.Train, error) {
	return slices.Collect(maps.Values(ts.trains)), nil
}

func (ts inMemTrainStore) GetTrainByRef(ref string) (exists bool, t shared.Train, err error) {
	t, exists = ts.trains[ref]
	return exists, t, nil
}

func (ts inMemTrainStore) RegisterListener(ref string, ctx context.Context) (chan shared.Train, error) {
	if strings.TrimSpace(ref) == "" {
		return nil, errorEmptyTrainRef
	}
	if _, ok := ts.trains[ref]; !ok {
		return nil, withTrainContext(errorNotFound)
	}

	sendCh := make(chan shared.Train)
	refListener := listener{
		sendCh: sendCh,
		ctx:    ctx,
	}

	listeners, ok := ts.listeners[ref]
	if !ok {
		ts.listeners[ref] = []listener{refListener}
	} else {
		listeners = append(listeners, refListener)
		ts.listeners[ref] = listeners
	}

	return sendCh, nil
}

func (ts inMemTrainStore) CreateTrain(t shared.Train) error {
	if strings.TrimSpace(t.Ref) == "" {
		return errorEmptyTrainRef
	}
	if _, exists := ts.trains[t.Ref]; exists {
		return withTrainContext(errorAlreadyExists)
	}
	ts.trains[t.Ref] = t
	return nil
}

func (ts inMemTrainStore) UpdateTrain(
	ref string,
	desc *string,
	posX, posY *int,
) (shared.Train, error) {
	if strings.TrimSpace(ref) == "" {
		return shared.Train{}, withTrainContext(errorEmptyTrainRef)
	}

	t, ok := ts.trains[ref]
	if !ok {
		return shared.Train{}, withTrainContext(errorNotFound)
	}

	if desc != nil {
		t.Description = *desc
	}

	if posX != nil {
		t.PosX = *posX
	}

	if posY != nil {
		t.PosY = *posY
	}

	listeners := ts.listeners[ref]

	listenerIndexesToRemove := []int{}
	for i, listener := range listeners {
		select {
		case <-listener.ctx.Done():
			listenerIndexesToRemove = append(listenerIndexesToRemove, i)
		default:
			listener.sendCh <- t
		}
	}

	for _, v := range listenerIndexesToRemove {
		log.Printf("Listener for train ref %v done, removing %v", ref, v)
		listeners = append(listeners[:v], listeners[v+1:]...)
		ts.listeners[ref] = listeners
	}

	ts.trains[ref] = t
	return t, nil
}
