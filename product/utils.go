package product

import (
	"github.com/zhuharev/micro/db"
)

func (s *Product) ReadRecentByType(t ProductType, limit uint64) (res []*Product, e error) {
	f := &db.Filter{Limit: limit, Order: "Id desc", WhereEq: db.WhereEq{}}
	f.WhereEq = append(f.WhereEq, map[string]interface{}{"Type": t})
	return s.ReadMany(f)
}
