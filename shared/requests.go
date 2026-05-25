package shared

import (
	"bytes"
	"encoding/json"
)

type (
	PatchTrainRequest struct {
		Description *string `json:"description"`
		PosX        *int    `json:"posX"`
		PosY        *int    `json:"posY"`
	}

	CreateStationRequest struct {
		Name string `json:"name"`
		PosX int    `json:"posX"`
		PosY int    `json:"posY"`
	}

	PatchStationRequest struct {
		Name string `json:"name"`
	}
)

func (r PatchTrainRequest) Read(b []byte) (int, error)    { return jsonRead(r, b) }
func (r CreateStationRequest) Read(b []byte) (int, error) { return jsonRead(r, b) }
func (r PatchStationRequest) Read(b []byte) (int, error)  { return jsonRead(r, b) }

func jsonRead(j any, b []byte) (int, error) {
	json, err := json.Marshal(j)
	if err != nil {
		return 0, err
	}
	return bytes.NewReader(json).Read(b)
}
