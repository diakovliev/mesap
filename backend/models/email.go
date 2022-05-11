package models

type Email struct {
	Id
	owner  IdData
	Active bool
	Mail   string
}
