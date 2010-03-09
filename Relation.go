package gouda

import (
        "container/vector"
//	"fmt"
	"reflect"
	"strings"
	)

type Relation struct {
	conditions vector.StringVector
	table string
}

func (r *Relation) Where(s string) *Relation {

	r.conditions.Push(s)

	return r

}

func (r *Relation) Table(t string) *Relation{
	r.table=t
return r
}

func (r *Relation) String() string {

	s := " Conditions : ("

	for _, ss := range r.conditions {
		s += ss
		if ss != r.conditions.Last() {
			s += ", "
		}
	}
	s += ")"

	return "Une relation" + s
}

func (r *Relation) Sql() (sql string) {
	sql = "Select * from "+r.table+" where ( "
	for _, ss := range r.conditions {
		sql += ss
		if ss != r.conditions.Last() {
			sql += " ) AND ( "
		}
	}

	sql += " );"

	return
}

func  NewRelation(t interface{}) (r *Relation){
	r=new(Relation)
	tab:=strings.Split(reflect.Typeof(t).String(),".",0)
	tablename:=strings.ToLower(tab[len(tab)-1])
	r.Table(tablename)
	return
}
