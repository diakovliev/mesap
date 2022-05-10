package models

type User struct {
	Id
	Login    string
	Salt     string
	Verifier string
}
