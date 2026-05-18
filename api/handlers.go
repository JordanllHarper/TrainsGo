package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/JordanllHarper/trainsgo/shared"
)

const getTrains = "GET /trains"

func handleGetTrains(tg trainGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trains, err := tg.GetTrains()
		if err != nil {
			internalServerError(w, err)
			return
		}
		writeJsonToHttpOk(w, trains)
	}
}

const getTrainByRef = "GET /trains/{ref}"

func handleGetTrainByRef(tg trainGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ref := r.PathValue("ref")
		exists, train, err := tg.GetTrainByRef(ref)
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

func handlePostTrain(tu trainCreater, srvBase string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t shared.Train
		if err := jsonDecode(r.Body, &t); err != nil {
			badRequest(w, err)
			return
		}

		switch err := tu.CreateTrain(t); {
		case errors.Is(err, errorEmptyRef):
			badRequest(w, err)
			return
		case errors.Is(err, errorAlreadyExists):
			http.Error(w, "Train with ref already exists", http.StatusConflict)
			return
		}

		writeJsonToHttpCreated(
			w,
			fmt.Sprintf("%v/trains/%v", srvBase, t.Ref),
			t,
		)
	}
}

const patchTrain = "PATCH /trains/{ref}"

func handlePatchTrain(tgu trainGetUpdater, sv secretVerifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		ref := r.PathValue("ref")

		var req shared.PatchTrainRequest
		if err := jsonDecode(r.Body, &req); err != nil {
			badRequest(w, err)
			return
		}

		exists, t, err := tgu.GetTrainByRef(ref)
		if err != nil {
			internalServerError(w, err)
			return
		}

		if !exists {
			http.NotFound(w, r)
			return
		}

		valid, err := sv.Verify(ref, auth)
		if err != nil {
			internalServerError(w, err)
			return
		}

		if !valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
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

		if _, err := tgu.UpsertTrain(t); err != nil {
			internalServerError(w, err)
			return
		}

		writeJsonToHttpOk(w, t)
	}
}
