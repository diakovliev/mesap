package models

type IdData = int64

const BAD_ID = -1
const FIRST_ID = 0

// Abstract Id
type Id struct {
	IdData `json:"Id"`
}

func (i *Id) SetId(id IdData) {
	i.IdData = id
}
func (i Id) GetId() IdData {
	return i.IdData
}
func MakeId(val IdData) Id {
	return Id{val}
}
