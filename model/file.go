package model

type File struct {
	ID        string `json:"id";sql:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name      string `json:"name"`
	Format    string `json:"format"`
	Size      int64  `json:"size"`
	Extension string `json:"extension"`
}
