package models

import "time"

type People struct {
	Id
	Name       string
	Surname    string
	Patronymic string
	Birth      time.Time
	Mail       string
}
