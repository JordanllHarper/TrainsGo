package shared

type PatchTrainRequest struct {
	Description *string `json:"description"`
	PosX        *int    `json:"posX"`
	PosY        *int    `json:"posY"`
}

type CreateStationRequest struct {
	Name string `json:"name"`
	PosX int    `json:"posX"`
	PosY int    `json:"posY"`
}

type PatchStationRequest struct {
	Name string `json:"name"`
}
