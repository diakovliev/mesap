package models

type Address struct {
	Id
	owner      IdData
	Active     bool
	Postcode   string
	Region     string
	City       string
	Street     string
	House      string
	Appartment string
}
