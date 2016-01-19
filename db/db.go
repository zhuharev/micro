package db

import (
	"errors"
)

type DB interface {
	Save(id int64, obj interface{}) error
	Read(int64, interface{}) (interface{}, error)

	ReadMany(interface{}, *Filter) ([]interface{}, error)

	Registry(interface{}, ...interface{}) error
}

type Filter struct {
	Offset uint64
	Limit  uint64
	Page   uint64

	Order   string
	WhereEq WhereEq
	Tags    []int64
}

type WhereEq []map[string]interface{}

func (w WhereEq) Append(name string, i interface{}) {
	w = append(w, map[string]interface{}{name: i})
}

type memoryDB map[int64]interface{}

func NewMemoryDb() memoryDB {
	return memoryDB{}
}

func (m memoryDB) Save(id int64, in interface{}, userId int64) error {
	m[id] = in
	return nil
}

func (m memoryDB) Open(id, userId int64) (interface{}, error) {
	obj, ok := m[id]
	if !ok {
		return nil, errors.New("Not found")
	}
	return obj, nil
}

func (m memoryDB) ReadMany(_ interface{}, f *Filter) (res []interface{}, e error) {
	for _, v := range m {
		res = append(res, v)
	}
	return res, nil
}
