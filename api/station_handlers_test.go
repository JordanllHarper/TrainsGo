package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/JordanllHarper/trainsgo/shared"
)

type mockStationGetter struct {
	stations            []shared.Station
	getStationById      shared.Station
	getStationExists    bool
	getStationError     error
	getStationByIdError error
}

func (m mockStationGetter) GetStations() ([]shared.Station, error) {
	return m.stations, m.getStationError
}
func (m mockStationGetter) GetStationById(ref string) (exists bool, t shared.Station, err error) {
	return m.getStationExists, m.getStationById, m.getStationByIdError
}

func TestHandleGetStations(t *testing.T) {
	populatedStations := []shared.Station{
		{
			Id:   "f5d2892a-d872-4520-84b0-6e20aae7c776",
			PosX: 0,
			PosY: 0,
		},
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		ts         stationGetter
		want       []shared.Station
		wantStatus int
		wantBody   bool
	}{
		{
			"Populated station store returns array of stations",
			mockStationGetter{stations: populatedStations},
			populatedStations,
			http.StatusOK,
			true,
		},
		{
			"Empty station store returns empty array",
			mockStationGetter{},
			[]shared.Station{},
			http.StatusOK,
			true,
		},
		{
			"Internal error returns internal server error",
			mockStationGetter{getStationError: errors.New("eek")},
			[]shared.Station{},
			http.StatusInternalServerError,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			handler := handleGetStations(tt.ts)

			req, err := http.NewRequest("GET", "/stations", nil)
			if err != nil {
				t.Fatal(err)
			}
			handler(recorder, req)
			if recorder.Code != tt.wantStatus {
				t.Errorf("HandleGetStations() status code = %v, want %v", recorder.Code, tt.wantStatus)
			}
			if tt.wantBody {
				var result []shared.Station
				if err = jsonDecode(recorder.Body, &result); err != nil {
					t.Fatal(err)
				}
				if !slices.Equal(result, tt.want) {
					t.Errorf("HandleGetStations() = %v, want %v", result, tt.want)
				}
			}
		})
	}
}

func TestHandleGetStationById(t *testing.T) {
	mockStation := shared.Station{
		Id:   "b211ce22-7310-434b-a453-2380521d6a7e",
		PosX: 0,
		PosY: 0,
	}
	tests := []struct {
		name            string
		ts              stationGetter
		wantStatusCode  int
		wantStationBody bool
		wantStation     shared.Station
	}{
		{
			"Populated station store returns matching station",
			mockStationGetter{
				getStationById:   mockStation,
				getStationExists: true,
			},
			http.StatusOK,
			true,
			mockStation,
		},
		{
			"Populated station store with no item matching ref returns not found",
			mockStationGetter{
				getStationById:   shared.Station{},
				getStationExists: false,
			},
			http.StatusNotFound,
			false,
			shared.Station{},
		},
		{
			"Internal error returns internal server error",
			mockStationGetter{
				getStationByIdError: errors.New("Eek"),
			},
			http.StatusInternalServerError,
			false,
			shared.Station{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			handler := handleGetStationById(tt.ts)

			req, err := http.NewRequest("GET", "/stations/{ref}", nil)
			req.SetPathValue("ref", mockStation.Id)
			if err != nil {
				t.Fatal(err)
			}
			handler(recorder, req)
			if recorder.Code != tt.wantStatusCode {
				t.Errorf("HandleGetStationById() status code = %v, want %v", recorder.Code, tt.wantStatusCode)
			}
			if tt.wantStationBody {
				var result shared.Station
				if err = jsonDecode(recorder.Body, &result); err != nil {
					t.Fatal(err)
				}
				if result != tt.wantStation {
					t.Errorf("HandleGetStationById() = %v, want %v", result, tt.wantStation)
				}
			}
		})
	}
}

type mockStationCreator struct {
	s   shared.Station
	err error
}

func (mtc mockStationCreator) CreateStation(name string, posX, posY int) (shared.Station, error) {
	return mtc.s, mtc.err
}

func TestHandlePostStation(t *testing.T) {
	s := shared.Station{
		Id:   "f5d2892a-d872-4520-84b0-6e20aae7c776",
		PosX: 0,
		PosY: 0,
	}

	tests := []struct {
		name            string
		ts              stationCreator
		wantStatusCode  int
		wantStationBody bool
		wantStation     shared.Station
		wantLocation    string
	}{
		{
			"Station already in x y returns conflict",
			mockStationCreator{err: errorAlreadyExists},
			http.StatusConflict,
			false,
			shared.Station{},
			"",
		},
		{
			"Station with same name returns conflict",
			mockStationCreator{err: errorAlreadyExists},
			http.StatusConflict,
			false,
			shared.Station{},
			"",
		},
		{
			"Doesnt exist returns 201 created",
			mockStationCreator{s: s},
			http.StatusCreated,
			true,
			s,
			"test/stations/f5d2892a-d872-4520-84b0-6e20aae7c776",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			handler := handlePostStation(tt.ts, "test")
			body, _ := json.Marshal(s)

			req, err := http.NewRequest("POST", "/stations", bytes.NewBuffer(body))
			req.SetPathValue("ref", "f5d2892a-d872-4520-84b0-6e20aae7c776")
			if err != nil {
				t.Fatal(err)
			}
			handler(recorder, req)

			if recorder.Code != tt.wantStatusCode {
				t.Errorf("HandlePostStation() status code = %v, want %v", recorder.Code, tt.wantStatusCode)
			}
			if tt.wantStationBody {
				var result shared.Station
				if err = jsonDecode(recorder.Body, &result); err != nil {
					t.Fatal(err)
				}
				if result != tt.wantStation {
					t.Errorf("HandlePostStation() = %v, want %v", result, tt.wantStation)
				}
			}
			location := recorder.Header().Get("location")
			if location != tt.wantLocation {
				t.Errorf("HandlePostStation() location = %v, want %v", location, tt.wantLocation)
			}
		})
	}
}

type mockStationUpdater struct {
	s   shared.Station
	err error
}

func (mtc mockStationUpdater) UpdateStation(id, name string) (shared.Station, error) {
	return mtc.s, mtc.err
}

func TestHandlePatchStation(t *testing.T) {
	renamedStation := shared.Station{
		Id:   "f5d2892a-d872-4520-84b0-6e20aae7c776",
		Name: "after",
		PosX: 0,
		PosY: 0,
	}

	patchReq := shared.PatchStationRequest{Name: "after"}
	badPatchReq := strings.NewReader("Bad json")

	tests := []struct {
		name            string
		patchReqBody    io.Reader
		ts              stationUpdater
		wantStatusCode  int
		wantStationBody bool
		wantStation     shared.Station
	}{
		{
			"Invalid json body returns bad request",
			badPatchReq,
			mockStationUpdater{},
			http.StatusBadRequest,
			false,
			shared.Station{},
		},
		{
			"Invalid id returns bad request",
			patchReq,
			mockStationUpdater{err: errorEmptyStationId},
			http.StatusBadRequest,
			false,
			shared.Station{},
		},
		{
			"Doesnt exist returns 404",
			patchReq,
			mockStationUpdater{err: errorNotFound},
			http.StatusNotFound,
			false,
			shared.Station{},
		},
		{
			"Internal error returns internal server error",
			patchReq,
			mockStationUpdater{err: errors.New("eek")},
			http.StatusInternalServerError,
			false,
			shared.Station{},
		},
		{
			"Station exists returns okay with new name",
			patchReq,
			mockStationUpdater{s: renamedStation},
			http.StatusOK,
			true,
			renamedStation,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			handler := handlePatchStationById(tt.ts)

			req, err := http.NewRequest("PATCH", "/stations/{id}", tt.patchReqBody)
			req.SetPathValue("id", "f5d2892a-d872-4520-84b0-6e20aae7c776")
			if err != nil {
				t.Fatal(err)
			}
			handler(recorder, req)

			if recorder.Code != tt.wantStatusCode {
				t.Errorf("handlePatchStationById() status code = %v, want %v", recorder.Code, tt.wantStatusCode)
			}
			if tt.wantStationBody {
				var result shared.Station
				if err = jsonDecode(recorder.Body, &result); err != nil {
					t.Fatal(err)
				}
				if result != tt.wantStation {
					t.Errorf("HandlePatchStationById() = %v, want %v", result, tt.wantStation)
				}
			}
		})
	}
}
