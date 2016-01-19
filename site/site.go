package site

import (
	"errors"
	"github.com/zhuharev/micro/db"
	"github.com/zhuharev/micro/obj"
	"github.com/zhuharev/micro/user"
	"log"
)

type Site struct {
	obj.Object

	Host string

	Ipv4 string
	Ipv6 string
}

func New(d db.DB) (*Site, error) {
	s, e := obj.NewInterface(d, &Site{})
	return s.(*Site), e
}

func (s *Site) Read(id int64) (*Site, error) {
	/*siteObj, e := s.DB.Read(id, *s)
	if e != nil {
		return nil, e
	}
	if site, ok := siteObj.(*Site); !ok {
		return nil, errors.New(ErrTypeConversion)
	} else {
		site.DB = s.DB
		site.Object.SelfPointer = site
		s = site
		return s, nil
	}
	return nil, errors.New("unknown")*/
	sObj, e := s.Object.Read(id)
	if e != nil {
		return nil, e
	}
	site := sObj.(*Site)
	//site.Object.SelfPointer = site.Object.SelfPointer.(*Site)
	return site, nil
}

type Service struct {
	db.DB
}

type WithUser struct {
	*Service
	User *user.User
}

func NewService(database db.DB) (*Service, error) {
	service := &Service{DB: database}
	return service, nil
}

func (srv *Service) New() (*Site, error) {
	return New(srv.DB)
}

func (srv *Service) WithUser(u *user.User) *WithUser {
	return &WithUser{Service: srv, User: u}
}

func (srv *WithUser) New() (*Site, error) {
	return New(srv.DB)
}

func (srv *WithUser) Read(id int64) (*Site, error) {
	if srv.DB == nil {
		return nil, errors.New(ErrDbIsNil)
	}
	o, err := srv.DB.Read(id, srv.User.Id)
	if err != nil {
		return nil, err
	}
	if o == nil {
		return nil, errors.New("Nil base result")
	}
	log.Println(o)
	site, ok := o.(*Site)
	if !ok {
		return nil, errors.New(ErrTypeConversion)
	}
	site.DB = srv.DB
	return site, nil
}

//func ()
