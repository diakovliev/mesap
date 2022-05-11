package models

import "time"

type GradeValue int
type PositionValue int

const (
	FirstGrade GradeValue = iota
	// -->
	Senjor
	Middle
	Junior
	// <-- last value
	LastGrade
)

func (g GradeValue) String() string {
	switch g {
	case Senjor:
		return "Senjor"
	case Middle:
		return "Middle"
	case Junior:
		return "Junior"
	}
	return UnknownValueString
}

func GradeFromString(input string) GradeValue {
	for i := FirstGrade + 1; i < LastGrade; i++ {
		if input == i.String() {
			return i
		}
	}
	return LastGrade
}

const (
	FirstPosition PositionValue = iota
	// -->
	Director
	Developer
	Ops
	Tester
	Office
	// <-- last value
	LastPosition
)

func (p PositionValue) String() string {
	switch p {
	case Director:
		return "Director"
	case Developer:
		return "Developer"
	case Ops:
		return "Ops"
	case Tester:
		return "Tester"
	case Office:
		return "Office"
	}
	return UnknownValueString
}

func PositionFromString(input string) PositionValue {
	for i := FirstPosition + 1; i < LastPosition; i++ {
		if input == i.String() {
			return i
		}
	}
	return LastPosition
}

type WorkingPeriod struct {
	Since *time.Time
	Till  *time.Time
}

type TaxInfo struct {
	RegisterDate time.Time
	Code         string
}

type People struct {
	Id

	// Basic
	Name       string
	Surname    string
	Patronymic string
	Birth      time.Time
	Photo      []byte // TODO:

	// References
	Phones       []Phone       // TODO: reference to dictionary
	Addresses    []Address     // TODO: reference to dictionary
	Emails       []Email       // TODO: reference to dictionary
	BankAccounts []BankAccount // TODO: reference to dictionary

	// Optionals ->
	Tax      []*TaxInfo       // TODO: reference to register
	Works    []*WorkingPeriod // TODO: reference to register
	Position *PositionValue   // TODO: reference to register
	Grade    *GradeValue      // TODO: reference to register
}
