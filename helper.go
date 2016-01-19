package micro

import (
	"github.com/zhuharev/micro/article"
	"github.com/zhuharev/micro/db"
	"github.com/zhuharev/micro/kv"
	"github.com/zhuharev/micro/product"
	"github.com/zhuharev/micro/user"
)

type Service struct {
	db db.DB
}

func NewService(db db.DB, objects ...interface{}) (*Service, error) {
	return &Service{db: db}, nil
}

func (s *Service) NewArticle(a ...article.Article) (*article.Article, error) {
	if len(a) > 0 {
		return article.NewFrom(a[0], s.db)
	}
	return article.New(s.db)
}

func (s *Service) NewProduct() (*product.Product, error) {
	return product.New(s.db)
}

func (s *Service) NewUser() (*user.User, error) {
	return user.New(s.db)
}

func (s *Service) NewKV() (*kv.KV, error) {
	return kv.New(s.db)
}
