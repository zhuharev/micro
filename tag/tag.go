package tag

import (
	"github.com/zhuharev/micro/obj"
)

type Tag struct {
	obj.Object
}

type Tags struct {
	obj.Object

	TagId    int64
	ItemType string //article, product, etc...
	ItemId   int64
}
