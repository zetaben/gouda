package gouda

import (
	"container/vector"
	"fmt"
	"reflect"
	"strings"
)

type RequestKind int
const (
SELECT RequestKind = iota
UPDATE
INSERT
COUNT
)

type Relation struct {
	conditions vector.StringVector
	table      string
	limit_offset int
	limit_count int
	order_field vector.StringVector
	order_direction vector.StringVector
	kind	RequestKind
	attributes vector.StringVector
}

func (r *Relation) Count(fields []string) *Relation {
	for _,s:=range fields {
	r.attributes.Push(s)
	}
	r.kind=COUNT
	return r
}

func (r *Relation) Where(s string) *Relation {

	r.conditions.Push(s)

	return r

}

func (r *Relation) Table(t string) *Relation {
	r.table = t
	return r
}

func (r * Relation ) Limit(offset, count int) *Relation {
	r.limit_offset=offset
	r.limit_count=count
	return r
}

func (r * Relation) First()  *Relation {
	r.Limit(0,1)
	return r
}

func (r *Relation) Order(field, direction string ) *Relation {
	r.order_field.Push(field)
	r.order_direction.Push(direction)
	return r
}

func (r *Relation) String() string {

	s := " Conditions : \n\t ("

	for _, ss := range r.conditions {
		s += ss
		if ss != r.conditions.Last() {
			s += ", "
		}
	}
	s += ") \n"
	if r.order_field.Len() > 0 {
		s += "Order :\n\t ("
		for i, ss := range r.order_field {
			s += ss
			s += " "
			s += r.order_direction[i]
			if ss != r.order_field.Last() {
				s += ", "
			}
		}
		s += ") \n"
	}
	if r.limit_count > 0 {
		s += "Offset : "+fmt.Sprint(r.limit_offset)+"\n"
		s += "Count : "+fmt.Sprint(r.limit_count)+"\n"
	}
	return "Une relation" + s
}

func From(t interface{}) (r *Relation) { return NewRelation(t) }

func NewRelation(t interface{}) (r *Relation) {
	r = new(Relation)
	r.kind=SELECT
	switch  typ := t.(type) {
	case string:
	r.Table(t.(string));
	default:
	tab := strings.Split(reflect.Typeof(t).String(), ".", 0)
	tablename := strings.ToLower(tab[len(tab)-1])
	r.Table(tablename)
	}
	return
}
