package model

type File struct {
	ID        string `json:"id";gorm:"primary_key"`
	Name      string `json:"name"`
	Format    string `json:"name"`
	Size      int32  `json:"size"`
	Extension string `json:"extension"`
}
