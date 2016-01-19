package article

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/zhuharev/micro/db"
	"github.com/zhuharev/micro/obj"
	//"html/template"
)

type Article struct {
	obj.Object

	PermanentUrl string
	Template     string

	Type ArticleType
	Body Body
}

type Body []byte

func (b Body) MarshalJSON() ([]byte, error) {
	var s = string(b)
	return json.MarshalIndent(&s, " ", " ")
}

func (u *Body) Scan(value interface{}) error {
	if value == nil {
		*u = Body{}
		return nil
	}
	*u = Body(value.([]byte))
	return nil
}
func (u Body) Value() (driver.Value, error) {
	return []byte(u), nil
}

//ArticleType also called as Category
type ArticleType int

const (
	NewsPost ArticleType = 1 << iota
	UserPost
	BlogPost
	CustomPage
)

func NewFrom(a Article, d db.DB) (*Article, error) {
	s, e := obj.NewInterface(d, &a)
	return s.(*Article), e
}

func New(d db.DB) (*Article, error) {
	a := Article{}
	return NewFrom(a, d)
}

func (s *Article) Read(id int64) (*Article, error) {
	sObj, e := s.Object.Read(id)
	if e != nil {
		return nil, e
	}
	site := sObj.(*Article)
	return site, nil
}

func (s *Article) ReadBySlug(id string) (*Article, error) {
	obj, e := s.Object.ReadBySlug(id)
	if e != nil {
		return nil, e
	}
	art := obj.(*Article)
	return art, nil
}

func (s *Article) ReadMany(f *db.Filter) (res []*Article, e error) {
	articles, e := s.Object.ReadMany(f)
	if e != nil {
		return nil, e
	}
	for i := range articles {
		res = append(res, articles[i].(*Article))
	}
	return
}

func (s *Article) ReadManyByType(t ArticleType, f *db.Filter) (res []*Article, e error) {
	if f == nil {
		f = &db.Filter{}
	}
	f.WhereEq = append(f.WhereEq, map[string]interface{}{"Type": t})

	return s.ReadMany(f)
}

//BodyStr useful for templates
func (s Article) BodyStr() string { //template.HTML {
	return string(s.Body) //template.HTML(string(s.Body))
}

func (s *Article) ReadByPermanentUrl(u string) (res *Article, e error) {

	f := &db.Filter{Limit: 1}

	f.WhereEq = append(f.WhereEq, map[string]interface{}{"PermanentUrl": u})

	r, e := s.ReadMany(f)
	if len(r) == 1 {
		return r[0], e
	}
	return nil, fmt.Errorf("Not found")
}

func (s *Article) ReadManyPage(page, pagesInPage uint64, filter *db.Filter) (res []*Article, e error) {
	articles, e := s.Object.ReadManyPage(page, pagesInPage, filter)
	if e != nil {
		return nil, e
	}
	for i := range articles {
		res = append(res, articles[i].(*Article))
	}
	return
}
