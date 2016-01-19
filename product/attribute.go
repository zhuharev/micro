package product

import (
	"encoding/json"
	"log"
	"sort"
)

type Attribute struct {
	Id    int64
	Name  string
	Value interface{}
	Type  AttributeType

	//Order weight
	Weight int
}

type AttributeType int

const (
	AString AttributeType = iota << 1
	AInt
)

func (s *Product) unmarshalAttributes() error {
	if s.attributes == nil {
		var v Attributes
		e := json.Unmarshal(s.Attributes, &v)
		if e != nil {
			return e
		}
		s.attributes = v
	}
	return nil
}

type Attributes []Attribute

func (a Attributes) Len() int           { return len(a) }
func (a Attributes) Less(i, j int) bool { return a[i].Weight < a[j].Weight }
func (a Attributes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (a Attributes) Attribute(name string) (Attribute, bool) {
	if a == nil {
		return Attribute{}, false
	}
	for _, v := range a {
		if v.Name == name {
			return v, true
		}
	}
	return Attribute{}, false
}

//Attribute useful for templates
func (s *Product) Attribute(name string) (Attribute, bool) { //template.HTML {
	e := s.unmarshalAttributes()
	if e != nil {
		log.Println(e)
		return Attribute{}, false
	}
	return s.attributes.Attribute(name)
}

func (s *Product) AddAttribute(name string, val interface{}, t AttributeType, weight int) {
	attrs := s.attributes
	attrs = append(attrs, Attribute{
		Name:  name,
		Value: val,
		Type:  t,

		//Order weight
		Weight: weight,
	})
	bts, e := json.Marshal(&attrs)
	if e != nil {
		panic(e)
		return
	}
	s.Attributes = bts
	s.attributes = attrs
}

//useful for templates
func (s *Product) MustAttribute(name string) Attribute {
	r, _ := s.Attribute(name)
	return r
}

func (s *Product) AllAttributes() Attributes { //template.HTML {
	e := s.unmarshalAttributes()
	if e != nil {
		return nil
	}
	sort.Sort(s.attributes)
	return s.attributes
}
