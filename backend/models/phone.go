package models

type PhoneType int

const (
	CityPhone PhoneType = iota
	MobilePhone
	Unknown
)

func (p PhoneType) String() string {
	switch p {
	case CityPhone:
		return "City"
	case MobilePhone:
		return "Mobile"
	case Unknown:
		return "Unknown"
	}
	return UnknownValueString
}

func PhoneTypeFromString(input string) PhoneType {
	switch input {
	case CityPhone.String():
		return CityPhone
	case MobilePhone.String():
		return MobilePhone
	}
	return Unknown
}

type Phone struct {
	Id
	owner  IdData
	Active bool
	Type   PhoneType
	Phone  string
}
