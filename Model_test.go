package gouda_test

import (
	"gouda"
	"testing"
	"fmt"
	"sort"
	"reflect"
)

type Personne struct {
	Nom string
	Id  int
}

func (p Personne) TableName() string { return "personne" }

func TestAttributes(t *testing.T) {
	p := new(Personne)
	attr := gouda.M(p).Attributes()
	types := map[string]reflect.Type{"Nom": &reflect.StringType{}, "Id": &reflect.IntType{}} //,"Age":&reflect.FloatType{}}

	if !compare(attr, types) {
		t.Error("Can't find Attributes found : " + fmt.Sprint(attr) + "\n")
	}

	var pp Personne
	attr = gouda.M(pp).Attributes()
	if !compare(attr, types) {
		t.Error("Can't find Attributes found : " + fmt.Sprint(attr) + "\n")
	}

}


func TestAttributesName(t *testing.T) {
	p := new(Personne)
	attr := gouda.M(p).AttributesNames()
	names := []string{"Nom","Id"}
	sort.StringArray(attr).Sort()
	sort.StringArray(names).Sort()
	if len(names)!=len(attr){
		  t.Error("Attributes Names, found  size mismatch : "+fmt.Sprint(len(attr))+" for "+fmt.Sprint(len(names)))
	}

	if !reflect.DeepEqual(names,attr) {
		  t.Error("Can't find Attributes Names, found : "+fmt.Sprint(attr))
	}
	var pp Personne
	attr = gouda.M(pp).AttributesNames()
	sort.StringArray(attr).Sort()
	if len(names)!=len(attr){
		  t.Error("Attributes Names, found  size mismatch : "+fmt.Sprint(len(attr))+" for "+fmt.Sprint(len(names)))
	}

	if !reflect.DeepEqual(names,attr) {
		  t.Error("Can't find Attributes Names, found : "+fmt.Sprint(attr))
	}
}


func TestModelName(t *testing.T) {
	p := new(Personne)
	if fname:=gouda.ModelName(p);fname!="Personne-personne" {
		  t.Error("wrong name found : "+fname)
	}

	var pp Personne
	if fname:=gouda.ModelName(pp);fname!="Personne-personne" {
		  t.Error("wrong name found : "+fname)
	}
}

func compare(a, b map[string]reflect.Type) bool {

	ok := true

	for k, v := range b {
		if _, present := a[k]; !present {
			ok = false
			break
		}

		if reflect.Typeof(a[k]) != reflect.Typeof(v) {
			ok = false
			break
		}
	}
	return ok
}
