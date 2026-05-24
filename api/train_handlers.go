package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/JordanllHarper/trainsgo/shared"
)

func setupTrainRoutes(
	mux *http.ServeMux,
	d dependencies,
	srvBase string,
) http.Handler {

	mux.HandleFunc(getTrains, handleGetTrains(d))
	mux.HandleFunc(getTrainByRef, handleGetTrainByRef(d))
	mux.HandleFunc(getTrainByRefLive, handleGetTrainByRefLive(d))
	mux.HandleFunc(postTrain, handlePostTrain(d, srvBase))
	mux.HandleFunc(patchTrain, addSecretValidation(handlePatchTrain(d, d)))

	return mux
}

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

const getTrainByRefLive = "GET /trains/{ref}/live"

func handleGetTrainByRefLive(tg trainUpdateSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ref := r.PathValue("ref")
		ch, err := tg.RegisterListener(ref, r.Context())
		switch {
		case errors.Is(err, errorEmptyTrainRef) || errors.Is(err, errorNotFound):
			badRequest(w, err)
			return
		case err != nil:
			internalServerError(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		clientGone := r.Context().Done()
		rc := http.NewResponseController(w)

		for {
			select {
			case <-clientGone:
				log.Println("Client disconnect")
				return
			case t := <-ch:
				if err := jsonEncode(w, t); err != nil {
					log.Println("error sending update to listeners for ref:", ref)
					return
				}
				if err := rc.Flush(); err != nil {
					log.Println("error pushing data to client:", err)
					return
				}
			}
		}
	}
}

const postTrain = "POST /trains"

func handlePostTrain(tu trainCreator, srvBase string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t shared.Train
		if err := jsonDecode(r.Body, &t); err != nil {
			badRequest(w, err)
			return
		}

		switch err := tu.CreateTrain(t); {
		case errors.Is(err, errorEmptyTrainRef):
			badRequest(w, err)
			return
		case errors.Is(err, errorAlreadyExists):
			http.Error(w, "Train with ref already exists", http.StatusConflict)
			return
		case err != nil:
			internalServerError(w, err)
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

func handlePatchTrain(tgu trainUpdater, sv secretVerifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		ref := r.PathValue("ref")

		var req shared.PatchTrainRequest
		if err := jsonDecode(r.Body, &req); err != nil {
			badRequest(w, err)
			return
		}

		valid, err := sv.Verify(ref, auth)
		if err != nil {
			internalServerError(w, err)
			return
		}

		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		t, err := tgu.UpdateTrain(
			ref,
			req.Description,
			req.PosX,
			req.PosY,
		)

		switch {
		case errors.Is(err, errorEmptyTrainRef):
			badRequest(w, err)
			return
		case err != nil:
			internalServerError(w, err)
			return
		}

		writeJsonToHttpOk(w, t)
	}
}
