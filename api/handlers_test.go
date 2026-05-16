package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/JordanllHarper/trainsgo/shared"
)

func TestHandleGetTrains(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		ts   trainStore
		want []shared.Train
	}{
		{
			"Populated train store returns array of trains",
			inMemTrainStore{
				"f5d2892a-d872-4520-84b0-6e20aae7c776": {
					Ref:         "f5d2892a-d872-4520-84b0-6e20aae7c776",
					Description: "test1",
					PosX:        0,
					PosY:        0,
				},
			},
			[]shared.Train{
				{
					Ref:         "f5d2892a-d872-4520-84b0-6e20aae7c776",
					Description: "test1",
					PosX:        0,
					PosY:        0,
				},
			},
		},

		{
			"Empty train store returns empty array",
			inMemTrainStore{},
			[]shared.Train{},
			// TODO: Add test cases.
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			handler := HandleGetTrains(tt.ts)

			req, err := http.NewRequest("GET", "/trains", nil)
			if err != nil {
				t.Fatal(err)
			}
			handler(recorder, req)
			var result []shared.Train
			if err = jsonDecode(recorder.Body, &result); err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(result, tt.want) {
				t.Errorf("HandleGetTrains() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestHandleGetTrainByRef(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		ts             trainStore
		wantStatusCode int
		wantTrainBody  bool
		wantTrain      shared.Train
	}{
		{
			"Populated train store returns matching train",
			inMemTrainStore{
				"f5d2892a-d872-4520-84b0-6e20aae7c776": {
					Ref:         "f5d2892a-d872-4520-84b0-6e20aae7c776",
					Description: "test1",
					PosX:        0,
					PosY:        0,
				},
			},
			http.StatusOK,
			true,
			shared.Train{
				Ref:         "f5d2892a-d872-4520-84b0-6e20aae7c776",
				Description: "test1",
				PosX:        0,
				PosY:        0,
			},
		},
		{
			"Empty train store returns not found",
			inMemTrainStore{},
			http.StatusNotFound,
			false,
			shared.Train{},
		},
		{
			"Populated train store with no item with ref returns not found",
			inMemTrainStore{
				"c54fa1a5-8cfc-46db-a083-8f5b9f1d3c09": {
					Ref:         "c54fa1a5-8cfc-46db-a083-8f5b9f1d3c09",
					Description: "test1",
					PosX:        0,
					PosY:        0,
				},
			},
			http.StatusNotFound,
			false,
			shared.Train{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			handler := HandleGetTrainByRef(tt.ts)

			req, err := http.NewRequest("GET", "/trains/{ref}", nil)
			req.SetPathValue("ref", "f5d2892a-d872-4520-84b0-6e20aae7c776")
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

func TestHandlePostTrain(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		ts             trainStore
		wantStatusCode int
		wantTrainBody  bool
		wantTrain      shared.Train
		wantLocation   string
	}{
		{
			"Populated train store returns conflict",
			inMemTrainStore{
				"f5d2892a-d872-4520-84b0-6e20aae7c776": {
					Ref:         "f5d2892a-d872-4520-84b0-6e20aae7c776",
					Description: "test1",
					PosX:        0,
					PosY:        0,
				},
			},
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
			handler := HandlePostTrain(tt.ts, "test")
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
