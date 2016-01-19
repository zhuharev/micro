package gorp

import (
	"github.com/stretchr/testify/assert"
	"github.com/zhuharev/micro/article"
	"github.com/zhuharev/micro/db"
	"github.com/zhuharev/micro/site"
	"testing"
)

var DB *GorpDb

func TestCreateTable(t *testing.T) {
	var e error
	DB, e = New("sqlite3", "db")
	if e != nil {
		t.Error(e)
	}
	e = DB.AddTable(site.Site{})
	if e != nil {
		t.Error(e)
	}

	DB.Registry(article.Article{})

}

func TestInsert(t *testing.T) {
	s, e := site.New(DB)
	if e != nil {
		t.Error(e)
	}
	s.Ipv4 = "1.1.1.1"
	e = s.Save()
	if e != nil {
		t.Error(e)
	}
	if s.Id == 0 {
		t.Error("auto incremet not working")
	}
}

func TestRead(t *testing.T) {
	s, e := site.New(DB)
	if e != nil {
		t.Error(e)
	}
	s1, e := s.Read(1)
	if e != nil {
		t.Error(e)
	}
	assert.Equal(t, s1.Id, int64(1), "Id must be 1")
}

func TestUpdate(t *testing.T) {
	s, e := site.New(DB)
	if e != nil {
		t.Error(e)
	}
	s.Id = 0
	s1, e := s.Read(1)
	if e != nil {
		t.Error(e)
	}
	//assert.Equal(t, s.Id, int64(1), "s id must be 1")
	assert.Equal(t, s1.Id, int64(1), "s1 id must be 1")
	assert.NotNil(t, s1.DB, "Db must be defined")
	assert.NotNil(t, s1.Object.SelfPointer, "SelfPointer must be defined")
	s1.Host = "google.com"
	t.Log(s1)
	e = s1.Save()
	if e != nil {
		t.Error(e)
	}

	s2, e := s1.Read(1)
	if e != nil {
		t.Error(e)
	}
	t.Log(s2)
}

func TestReadMany(t *testing.T) {
	art, e := article.New(DB)
	if e != nil {
		t.Error(e)
	}
	art.Title = "Art 1"
	e = art.Save()
	if e != nil {
		t.Error(e)
	}

	art1, _ := article.New(DB)

	art1.Title = "Art 2"
	e = art1.Save()
	if e != nil {
		t.Error(e)
	}

	arts, e := art.ReadMany(&db.Filter{Limit: 2})
	if e != nil {
		t.Error(e)
	}

	assert.Equal(t, 2, len(arts), "Sould be 2 articles")
	art = arts[0]

	art.Title = "Kirya"
	e = art.Save()
	if e != nil {
		t.Error(e)
	}
}
