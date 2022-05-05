package ifaces

import "github.com/diakovliev/mesap/backend/models"

type Id interface {
	SetId(models.IdData)
	GetId() models.IdData
}

type Models interface {
	models.User | models.People | models.Role
}

type Table[M Models] interface {
	Get(models.IdData) (M, error)
	Find(func(record M) bool) (M, error)
	Each(func(record M) bool) error
	Insert(record M) (models.IdData, error)
	Update(record M) error
	Delete(models.IdData) error
}

type Database interface {
	Open() error
	Close()

	Users() (Table[models.User], error)
	Peoples() (Table[models.People], error)
	Roles() (Table[models.Role], error)
}
