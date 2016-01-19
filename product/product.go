package product

import (
	"github.com/zhuharev/micro/article"
	"github.com/zhuharev/micro/db"
	"github.com/zhuharev/micro/obj"
	//"html/template"
)

type Product struct {
	obj.Object

	Type ProductType

	Attributes []byte
	Body       article.Body
	attributes Attributes `db:"-"`
}

type ProductType int

func New(d db.DB) (*Product, error) {
	s, e := obj.NewInterface(d, &Product{})
	return s.(*Product), e
}

func (s *Product) Read(id int64) (*Product, error) {
	sObj, e := s.Object.Read(id)
	if e != nil {
		return nil, e
	}
	site := sObj.(*Product)
	return site, nil
}

func (s *Product) ReadBySlug(id string) (*Product, error) {
	obj, e := s.Object.ReadBySlug(id)
	if e != nil {
		return nil, e
	}
	art := obj.(*Product)
	return art, nil
}

func (s *Product) ReadMany(f *db.Filter) (res []*Product, e error) {
	Products, e := s.Object.ReadMany(f)
	if e != nil {
		return nil, e
	}
	for i := range Products {
		res = append(res, Products[i].(*Product))
	}
	return
}

func (s *Product) ReadManyByType(t ProductType, f *db.Filter) (res []*Product, e error) {
	if f == nil {
		f = &db.Filter{}
	}
	f.WhereEq = append(f.WhereEq, map[string]interface{}{"Type": t})

	return s.ReadMany(f)
}
