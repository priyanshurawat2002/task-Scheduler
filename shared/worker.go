package shared

type Worker struct {
	ID    string `json:"id"`
	Load  int    `json:"load"`
	Alive bool   `json:"alive"`
}
