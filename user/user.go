package user

import (
	"github.com/zhuharev/micro/db"
	"github.com/zhuharev/micro/obj"
)

type User struct {
	obj.Object

	Username  string
	FirstName string

	HashedPassword []byte
}

func New(d db.DB) (*User, error) {
	s, e := obj.NewInterface(d, &User{})
	return s.(*User), e
}

func (u *User) Read(id int64) (*User, error) {
	obj, e := u.Object.Read(id)
	if e != nil {
		return nil, e
	}
	user := obj.(*User)
	return user, nil
}

func (u *User) ReadMany(f *db.Filter) (res []*User, e error) {
	users, e := u.Object.ReadMany(f)
	if e != nil {
		return nil, e
	}
	for i := range users {
		res = append(res, users[i].(*User))
	}
	return
}
