package main

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	badRef = errors.New("Bad ref")
)

const getTrains = "GET /trains"

func HandleGetTrains(ts trainGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trains, err := ts.GetTrains()
		if err != nil {
			badRequest(w, err)
			return
		}
		writeJsonToHttpOk(w, trains)
	}
}

const getTrainByRef = "GET /trains/{ref}"

func HandleGetTrainByRef(ts trainGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ref := r.PathValue("ref")
		if ref == "" {
			badRequest(w, badRef)
			return
		}

		exists, train, err := ts.GetTrainByRef(ref)
		if err != nil {
			internalServerError(w, err)
			return
		}

		if !exists {
			http.NotFound(w, r)
			return
		}

		writeJsonToHttpOk(w, train)
	}
}

const postTrain = "POST /trains"

func HandlePostTrain(ts trainGetUpdater, srvBase string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t Train
		if err := jsonDecode(r.Body, &t); err != nil {
			badRequest(w, err)
			return
		}

		exists, _, err := ts.GetTrainByRef(t.Ref)
		if err != nil {
			internalServerError(w, err)
			return
		}

		if exists {
			http.Error(w, "Train with ref already exists", http.StatusConflict)
			return
		}

		if _, err = ts.UpsertTrain(t); err != nil {
			internalServerError(w, err)
			return
		}

		location := fmt.Sprintf("%v/trains/%v", srvBase, t.Ref)

		writeJsonToHttpCreated(w, location, t)
	}
}

type PatchTrainRequest struct {
	Description *string `json:"description"`
	PosX        *int    `json:"posX"`
	PosY        *int    `json:"posY"`
}

const patchTrain = "PATCH /trains/{ref}"

func HandlePatchTrain(ts trainGetUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ref := r.PathValue("ref")
		if ref == "" {
			badRequest(w, badRef)
			return
		}

		var req PatchTrainRequest
		if err := jsonDecode(r.Body, req); err != nil {
			badRequest(w, err)
			return
		}

		exists, t, err := ts.GetTrainByRef(ref)
		if err != nil {
			badRequest(w, badRef)
			return
		}

		if !exists {
			http.NotFound(w, r)
			return
		}

		if req.Description != nil {
			t.Description = *req.Description
		}

		if req.PosX != nil {
			t.PosX = *req.PosX
		}

		if req.PosY != nil {
			t.PosY = *req.PosY
		}

		if _, err := ts.UpsertTrain(t); err != nil {
			internalServerError(w, err)
			return
		}

		writeJsonToHttpOk(w, t)
	}
}
