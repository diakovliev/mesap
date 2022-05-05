package fake_database

import (
	"errors"
	"fmt"
	"sync"

	"github.com/diakovliev/mesap/backend/ifaces"
	"github.com/diakovliev/mesap/backend/models"
)

var (
	ErrWrongRecord  = errors.New("Wrong record!")
	ErrEmptyTable   = errors.New("Table is empty!")
	ErrNoSuchRecord = errors.New("No such record!")
)

type FakeTable[M ifaces.Models] struct {
	sync.Mutex
	parent sync.Locker
	table  map[models.IdData]*M
	currId models.IdData
}

func makeFakeTable[M ifaces.Models](initialId models.IdData) FakeTable[M] {
	return FakeTable[M]{
		parent: nil,
		table:  make(map[models.IdData]*M),
		currId: initialId,
	}
}

func NewFakeTable[M ifaces.Models](initialId models.IdData) *FakeTable[M] {
	res := makeFakeTable[M](initialId)
	return &res
}

func (T *FakeTable[M]) nextId() models.IdData {
	ret := T.currId
	T.currId += 1
	return ret
}

func (T *FakeTable[M]) Insert(record M) (models.IdData, error) {
	T.parent.Lock()
	T.Mutex.Lock()
	defer func() {
		T.Mutex.Unlock()
		T.parent.Unlock()
	}()

	var i interface{} = &record

	id, ok := i.(ifaces.Id)
	if !ok {
		return models.BAD_ID, ErrWrongRecord
	}

	id.SetId(T.nextId())

	_, ok = T.table[id.GetId()]
	if ok {
		panic(fmt.Errorf("Record with id: %d already exist!", id.GetId()))
	}

	T.table[id.GetId()] = &record

	return id.GetId(), nil
}

func (T *FakeTable[M]) Update(record M) error {
	T.parent.Lock()
	T.Mutex.Lock()
	defer func() {
		T.Mutex.Unlock()
		T.parent.Unlock()
	}()

	var i interface{} = &record

	id, ok := i.(ifaces.Id)
	if !ok {
		return ErrWrongRecord
	}

	_, ok = T.table[id.GetId()]
	if !ok {
		return fmt.Errorf("Record with id: %d not exist!", id.GetId())
	}

	T.table[id.GetId()] = &record
	return nil
}

func (T *FakeTable[M]) Delete(id models.IdData) error {
	T.parent.Lock()
	T.Mutex.Lock()
	defer func() {
		T.Mutex.Unlock()
		T.parent.Unlock()
	}()

	_, ok := T.table[id]
	if !ok {
		return ErrNoSuchRecord
	}

	delete(T.table, id)

	return nil
}

func (T *FakeTable[M]) Get(id models.IdData) (M, error) {
	T.parent.Lock()
	T.Mutex.Lock()
	defer func() {
		T.Mutex.Unlock()
		T.parent.Unlock()
	}()

	var res M

	ret, ok := T.table[id]
	if !ok {
		return res, ErrNoSuchRecord
	}

	res = *ret

	return res, nil
}

func (T *FakeTable[M]) Each(callback func(record M) bool) error {
	T.parent.Lock()
	T.Mutex.Lock()
	defer func() {
		T.Mutex.Unlock()
		T.parent.Unlock()
	}()

	err := ErrEmptyTable

	for _, record := range T.table {
		err = nil
		if !callback(*record) {
			break
		}
	}

	return err
}

func (T *FakeTable[M]) Find(callback func(record M) bool) (M, error) {
	T.parent.Lock()
	T.Mutex.Lock()
	defer func() {
		T.Mutex.Unlock()
		T.parent.Unlock()
	}()

	err := ErrEmptyTable

	var res M
	for _, record := range T.table {
		err = nil

		if callback(*record) {
			res = *record
			break
		}
	}

	return res, err
}

///////////////////////////////////////////////////////////////////////////////
type FakeUsers struct {
	FakeTable[models.User]
}
type FakePeoples struct {
	FakeTable[models.People]
}
type FakeRoles struct {
	FakeTable[models.Role]
}

type FakeDatabase struct {
	sync.Mutex
	users   *FakeUsers
	peoples *FakePeoples
	roles   *FakeRoles
}

///////////////////////////////////////////////////////////////////////////////
func NewDatabase() ifaces.Database {
	ret := &FakeDatabase{
		users:   &FakeUsers{FakeTable: makeFakeTable[models.User](models.FIRST_ID)},
		peoples: &FakePeoples{FakeTable: makeFakeTable[models.People](models.FIRST_ID)},
		roles:   &FakeRoles{FakeTable: makeFakeTable[models.Role](models.FIRST_ID)},
	}
	ret.users.parent = ret
	ret.peoples.parent = ret
	ret.roles.parent = ret
	return ret
}
func (*FakeDatabase) Open() error {
	return nil
}
func (*FakeDatabase) Close() {
}
func (d *FakeDatabase) Users() (ifaces.Table[models.User], error) {
	return d.users, nil
}
func (d *FakeDatabase) Peoples() (ifaces.Table[models.People], error) {
	return d.peoples, nil
}
func (d *FakeDatabase) Roles() (ifaces.Table[models.Role], error) {
	return d.roles, nil
}
