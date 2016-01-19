package obj

import (
	"errors"
	"fmt"
	"github.com/sisteamnik/guseful/chpu"
	"github.com/zhuharev/micro/db"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	ErrDbIsNil        = "Database is nil"
	ErrTypeConversion = "Cannot converse type"
)

type Object struct {
	Id          int64
	Title       string
	Description string
	Image       int64

	//organization or private person if 0 -> nobody
	Owner int64

	ObjectDataMeta

	db.DB       `json:"-" db:"-"`
	SelfPointer `json:"-" db:"-"`
}

type SelfPointer interface{}

//field used by database
type ObjectDataMeta struct {
	RecordCreated   int64
	RecordModified  int64
	RecordPublished int64

	RecordCreatedByUser int64
	RecordAclId         int64
}

func New(d db.DB) Object {
	t := time.Now()
	un := t.UnixNano()
	o := Object{}
	o.RecordCreated = un
	o.RecordModified = un

	o.DB = d
	return o
}

//code from https://github.com/brendensoares/storm/blob/master/model.go
func NewInterface(d db.DB, model interface{}) (interface{}, error) {
	modelValue := reflect.ValueOf(model)
	elemModel := modelValue.Elem() //modelValue.Elem() // reflect.New(modelValue.Type()).Elem()

	o := New(d)
	o.SelfPointer = model
	obj := reflect.ValueOf(o)

	objField := elemModel.FieldByName("Object")
	if !objField.CanSet() {
		panic("cant set")
	} else {
		objField.Set(obj)
	}
	if !elemModel.CanAddr() {
		panic("cant addr")
	}
	ptr := elemModel.Addr().Interface()
	return ptr, nil
}

func (o *Object) Save() error {
	return o.DB.Save(o.Id, o.SelfPointer)
}

func (o *Object) Read(id int64) (interface{}, error) {
	if o.SelfPointer != nil {
		db := o.DB
		var e error
		vSP := reflect.ValueOf(o.SelfPointer).Elem()
		//vPtrSP := reflect.New(vSP.Type())
		newvSPIf, e := o.DB.Read(id, vSP.Interface())
		if e != nil {
			return nil, e
		}

		assign(newvSPIf, db)
		return newvSPIf, e
	}
	log.Print()
	return nil, errors.New("SelfPointer is nil")
}

func (o Object) Slug() string {
	return chpu.Chpu(o.Title) + "_" + fmt.Sprint(o.Id)
}

func (o *Object) ReadBySlug(slug string) (interface{}, error) {
	getId := func(s string) int64 {
		arr := strings.Split(s, "_")
		strId := arr[len(arr)-1]
		//todo handle error
		id, _ := strconv.ParseInt(strId, 10, 64)
		return id
	}
	id := getId(slug)
	if id == 0 {
		return nil, errors.New("Not found")
	}
	return o.Read(id)
}

func assign(i interface{}, d db.DB) {
	val := reflect.ValueOf(i).Elem()
	val.FieldByName("DB").Set(reflect.ValueOf(d))
	field := val.FieldByName("SelfPointer")
	field.Set(reflect.ValueOf(i))
}

func (o *Object) ReadMany(f *db.Filter) ([]interface{}, error) {
	if o.SelfPointer == nil {
		return nil, errors.New("SelfPointer is nil")
	}
	spValue := reflect.ValueOf(o.SelfPointer).Elem()
	//spType := reflect.TypeOf(o.SelfPointer)
	//reflect.Indirect(reflect.ValueOf(i))
	//sliceValue := reflect.MakeSlice(reflect.SliceOf(spType), 0, 0)
	//el := reflect.ValueOf(o.SelfPointer).Elem()
	els, e := o.DB.ReadMany(spValue.Interface(), f)
	if e != nil {
		return nil, e
	}
	for i := range els {
		assign(els[i], o.DB)
	}
	return els, nil
}

/*
objField := elemModel.FieldByName("Object")
	if !objField.CanSet() {
		panic("cant set")
	} else {
		objField.Set(obj)
	}
	if !elemModel.CanAddr() {
		panic("cant addr")
	}
	ptr := elemModel.Addr().Interface()*/
