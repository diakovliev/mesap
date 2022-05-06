package models

type User struct {
	Id
	Login string
	Mail  string
	Ii    string // SRP-6a identifier
	Iv    string // SRP-6a verifier
}
