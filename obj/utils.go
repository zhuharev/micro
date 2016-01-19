package obj

import (
	"github.com/zhuharev/micro/db"
)

func (o *Object) ReadManyPage(page, pagesInPage uint64, filter *db.Filter) ([]interface{}, error) {
	if filter == nil {
		filter = &db.Filter{}
	}

	filter.Limit = pagesInPage
	filter.Offset = (page - 1) * pagesInPage

	return o.ReadMany(filter)
}
