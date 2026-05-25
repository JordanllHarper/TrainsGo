package shared

type Train struct {
	Ref         string `json:"ref"`
	Description string `json:"description"`
	PosX        int    `json:"posX"`
	PosY        int    `json:"posY"`
}

type Station struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	PosX int    `json:"posX"`
	PosY int    `json:"posY"`
}

func (t Train) Read(b []byte) (int, error) { return jsonRead(t, b) }
