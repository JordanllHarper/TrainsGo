package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/JordanllHarper/trainsgo/shared"
)

type mockTrainGetter struct {
	trains             []shared.Train
	getTrainByRef      shared.Train
	getTrainExists     bool
	getTrainError      error
	getTrainByRefError error
}

func (m mockTrainGetter) GetTrains() ([]shared.Train, error) {
	return m.trains, m.getTrainError
}
func (m mockTrainGetter) GetTrainByRef(ref string) (exists bool, t shared.Train, err error) {
	return m.getTrainExists, m.getTrainByRef, m.getTrainByRefError
}

func TestHandleGetTrains(t *testing.T) {
	populatedTrains := []shared.Train{
		{
			Ref:         "f5d2892a-d872-4520-84b0-6e20aae7c776",
			Description: "test1",
			PosX:        0,
			PosY:        0,
		},
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		ts         trainGetter
		want       []shared.Train
		wantStatus int
		wantBody   bool
	}{
		{
			"Populated train store returns array of trains",
			mockTrainGetter{trains: populatedTrains},
			populatedTrains,
			http.StatusOK,
			true,
		},
		{
			"Empty train store returns empty array",
			mockTrainGetter{},
			[]shared.Train{},
			http.StatusOK,
			true,
		},
		{
			"Internal error returns internal server error",
			mockTrainGetter{getTrainError: errors.New("eek")},
			[]shared.Train{},
			http.StatusInternalServerError,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			handler := handleGetTrains(tt.ts)

			req, err := http.NewRequest("GET", "/trains", nil)
			if err != nil {
				t.Fatal(err)
			}
			handler(recorder, req)
			if recorder.Code != tt.wantStatus {
				t.Errorf("HandleGetTrains() status code = %v, want %v", recorder.Code, tt.wantStatus)
			}
			if tt.wantBody {
				var result []shared.Train
				if err = jsonDecode(recorder.Body, &result); err != nil {
					t.Fatal(err)
				}
				if !slices.Equal(result, tt.want) {
					t.Errorf("HandleGetTrains() = %v, want %v", result, tt.want)
				}
			}
		})
	}
}

func TestHandleGetTrainByRef(t *testing.T) {
	mockTrain := shared.Train{
		Ref:         "f5d2892a-d872-4520-84b0-6e20aae7c776",
		Description: "test1",
		PosX:        0,
		PosY:        0,
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		ts             trainGetter
		wantStatusCode int
		wantTrainBody  bool
		wantTrain      shared.Train
	}{
		{
			"Populated train store returns matching train",
			mockTrainGetter{
				getTrainByRef:  mockTrain,
				getTrainExists: true,
			},
			http.StatusOK,
			true,
			mockTrain,
		},
		{
			"Populated train store with no item matching ref returns not found",
			mockTrainGetter{
				getTrainByRef:  shared.Train{},
				getTrainExists: false,
			},
			http.StatusNotFound,
			false,
			shared.Train{},
		},
		{
			"Internal error returns internal server error",
			mockTrainGetter{
				getTrainByRefError: errors.New("Eek"),
			},
			http.StatusInternalServerError,
			false,
			shared.Train{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			handler := handleGetTrainByRef(tt.ts)

			req, err := http.NewRequest("GET", "/trains/{ref}", nil)
			req.SetPathValue("ref", mockTrain.Ref)
			if err != nil {
				t.Fatal(err)
			}
			handler(recorder, req)
			if recorder.Code != tt.wantStatusCode {
				t.Errorf("HandleGetTrainByRef() status code = %v, want %v", recorder.Code, tt.wantStatusCode)
			}
			if tt.wantTrainBody {
				var result shared.Train
				if err = jsonDecode(recorder.Body, &result); err != nil {
					t.Fatal(err)
				}
				if result != tt.wantTrain {
					t.Errorf("HandleGetTrainByRef() = %v, want %v", result, tt.wantTrain)
				}
			}
		})
	}
}

type mockTrainCreator struct {
	err error
}

func (mtc mockTrainCreator) CreateTrain(t shared.Train) error { return mtc.err }

func TestHandlePostTrain(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		ts             trainCreater
		wantStatusCode int
		wantTrainBody  bool
		wantTrain      shared.Train
		wantLocation   string
	}{
		{
			"Populated train store returns conflict",
			mockTrainCreator{errorAlreadyExists},
			http.StatusConflict,
			false,
			shared.Train{},
			"",
		},
		{
			"Empty train store return 201 created",
			inMemTrainStore{},
			http.StatusCreated,
			true,
			shared.Train{
				Ref:         "f5d2892a-d872-4520-84b0-6e20aae7c776",
				Description: "test1",
				PosX:        0,
				PosY:        0,
			},
			"test/trains/f5d2892a-d872-4520-84b0-6e20aae7c776",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			handler := handlePostTrain(tt.ts, "test")
			train := shared.Train{
				Ref:         "f5d2892a-d872-4520-84b0-6e20aae7c776",
				Description: "test1",
				PosX:        0,
				PosY:        0,
			}

			body, _ := json.Marshal(train)

			req, err := http.NewRequest("POST", "/trains", bytes.NewBuffer(body))
			req.SetPathValue("ref", "f5d2892a-d872-4520-84b0-6e20aae7c776")
			if err != nil {
				t.Fatal(err)
			}
			handler(recorder, req)

			if recorder.Code != tt.wantStatusCode {
				t.Errorf("HandlePostTrain() status code = %v, want %v", recorder.Code, tt.wantStatusCode)
			}
			if tt.wantTrainBody {
				var result shared.Train
				if err = jsonDecode(recorder.Body, &result); err != nil {
					t.Fatal(err)
				}
				if result != tt.wantTrain {
					t.Errorf("HandlePostTrain() = %v, want %v", result, tt.wantTrain)
				}
			}
			location := recorder.Header().Get("location")
			if location != tt.wantLocation {
				t.Errorf("HandlePostTrain() location = %v, want %v", location, tt.wantLocation)
			}
		})
	}
}
