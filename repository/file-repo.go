package repository

import "github.com/zeyd17/file-microservice/model"

type IFileRepo interface {
	GetByID(id string) (*model.File, error)
	Create(f *model.File) error
	Delete(id string) (bool, error)
}
