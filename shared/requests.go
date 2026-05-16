package shared

type PatchTrainRequest struct {
	Description *string `json:"description"`
	PosX        *int    `json:"posX"`
	PosY        *int    `json:"posY"`
}
