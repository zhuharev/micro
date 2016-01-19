package acl

import (
	"github.com/zhuharev/micro/obj"
)

type AccessControlList interface {
	CanWrite(in interface{}, userId int64) bool
	CanRead(in interface{}, userId int64) bool
}

type Acl struct {
	obj.Object

	AllowWriteAccess []int64
	AllowReadAccess  []int64
	DenyWriteAccess  []int64
	DenyReadAccess   []int64
}

func (a *Acl) CanWrite(in interface{}, userId int64) {
	if !a.CanRead(in, userId) {
		return false
	}
	if len(a.AllowWriteAccess) == 0 {
		if !in(userId, a.DenyWriteAccess) {
			return true
		}
	}
	return false
}

func (a *Acl) CanRead(in interface{}, userId int64) {
	if len(a.AllowReadAccess) == 0 {
		if !in(userId, a.DenyReadAccess) {
			return true
		}
	}
	return false
}

func in(id int64, arr []int64) bool {
	for _, v := range arr {
		if v == id {
			return true
		}
	}
	return false
}
