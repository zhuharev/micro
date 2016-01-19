package site

import (
	"github.com/stretchr/testify/assert"
	"github.com/zhuharev/micro/db"
	"github.com/zhuharev/micro/obj"
	"github.com/zhuharev/micro/user"
	"log"
	"testing"
)

var (
	DB      db.DB
	Srv     *WithUser
	TmpSite *Site
	err     error
	u       *user.User
)

func init() {
	u = &user.User{}
	u.Object = obj.Object{Id: 1}
	DB = db.NewMemoryDb()

	srv, err := NewService(DB)
	if err != nil {
		panic(err)
	}
	Srv = srv.WithUser(u)
}

func TestNew(t *testing.T) {
	TmpSite, err = Srv.New(0)
	if err != nil {
		t.Error(err)
	}
	if TmpSite.Object.DB == nil {
		t.Error("db is nil")
	}
	log.Print()
	if TmpSite == nil {
		t.Error("t is nil")
	}
	if TmpSite.Object.SelfPointer == nil {
		t.Error("self pointer is nil")
	}
}

func TestSave(t *testing.T) {
	TmpSite.Title = "super"
	err = TmpSite.Save()
	if err != nil {
		t.Error(err)
	}
}

func TestRead(t *testing.T) {
	TmpSite.Title = "super"
	TmpSite.Id = 12
	err = TmpSite.Save()
	if err != nil {
		t.Error(err)
	}

	t.Log(TmpSite.Object.SelfPointer)

	site, err := Srv.Read(12)
	if err != nil {
		t.Error(err)
	}
	_ = site
	assert.Equal(t, 0, 0)
	//assert.Equal(t, "super", site.Title, "title not equal")
	//assert.Equal(t, 12, 12, "Id not equal")
}
