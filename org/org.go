package org

import (
	"github.com/zhuharev/micro/obj"
)

type Organization struct {
	obj.Object

	//individual person or entity
	Type string
}
