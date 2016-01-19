package kv

import (
	"encoding/json"
	"fmt"
	"github.com/zhuharev/micro/db"
	"github.com/zhuharev/micro/obj"
)

type KV struct {
	obj.Object
	Key   string
	Value Value

	Type KVType
}

type Value []byte

func (v Value) MarshalJSON() ([]byte, error) {
	s := string(v)
	return json.Marshal(&s)
}

type KVType int

const (
	KVSettings KVType = iota << 1
	KVData
)

func New(d db.DB) (*KV, error) {
	s, e := obj.NewInterface(d, &KV{})
	return s.(*KV), e
}

func (s *KV) Read(id int64) (*KV, error) {
	sObj, e := s.Object.Read(id)
	if e != nil {
		return nil, e
	}
	site := sObj.(*KV)
	return site, nil
}

func (s *KV) ReadByKey(key string) (*KV, error) {
	f := &db.Filter{}
	f.WhereEq = append(f.WhereEq, map[string]interface{}{"Key": key})
	f.Limit = 1
	arr, e := s.ReadMany(f)
	if len(arr) == 1 {
		return arr[0], e
	}
	return nil, fmt.Errorf("Not found")
}

func (s *KV) ReadMany(f *db.Filter) (res []*KV, e error) {
	KVs, e := s.Object.ReadMany(f)
	if e != nil {
		return nil, e
	}
	for i := range KVs {
		res = append(res, KVs[i].(*KV))
	}
	return
}

func (s *KV) ReadManyByType(t KVType, f *db.Filter) (res []*KV, e error) {
	if f == nil {
		f = &db.Filter{}
	}
	f.WhereEq = append(f.WhereEq, map[string]interface{}{"Type": t})

	return s.ReadMany(f)
}
