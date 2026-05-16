package shared

type Train struct {
	Ref         string `json:"ref"`
	Description string `json:"description"`
	PosX        int    `json:"posX"`
	PosY        int    `json:"posY"`
}
