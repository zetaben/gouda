package gouda_test

import (
	"gouda"
	"testing"
	"fmt"
	"os"
	"reflect"
)

var testFileName string = "/tmp/test.xml"

func TestOpenClose(t *testing.T) {
	conn := gouda.OpenXML(testFileName)
	conn.Close()
	if _, err := os.Open(testFileName, os.O_RDONLY, 0); err != nil {
		t.Error("database file not created !")
	} else {
		os.Remove(testFileName)
	}
}


func TestCreateTable(t *testing.T) {
	conn := gouda.OpenXML(testFileName).(*gouda.XMLConnector)
	table := conn.CreateTable("personne")
	table.AddAttribute("nom", gouda.StringKind)
	table.AddAttribute("age", gouda.IntKind)
	//	fmt.Println(table)
	conn.Close()
	conn = gouda.OpenXML(testFileName).(*gouda.XMLConnector)
	if conn.Table("personne") == nil {
		t.Error("Can't find Table")
	}

	if attr := conn.Table("personne").Attributes(); !reflect.DeepEqual(attr, table.Attributes()) {
		t.Error("not found attrs " + fmt.Sprint(attr) + " vs " + fmt.Sprint(table.Attributes()))
	}
	conn.Close()
}

func TestCreateData(t *testing.T) {
	conn := gouda.OpenXML(testFileName).(*gouda.XMLConnector)
	table := conn.CreateTable("personne")
	table.AddAttribute("nom", gouda.StringKind)
	table.AddAttribute("age", gouda.IntKind)
	//	fmt.Println(table)
	table.Insert(map[string]gouda.Value{"id": gouda.SysInt(1).Value(), "nom": gouda.SysString("toto").Value(), "age": gouda.SysInt(13).Value()})
	table.Insert(map[string]gouda.Value{"id": gouda.SysInt(2).Value(), "nom": gouda.SysString("titi").Value(), "age": gouda.SysInt(13).Value()})
	conn.Close()
	conn = gouda.OpenXML(testFileName).(*gouda.XMLConnector)
	if conn.Table("personne").Data().Len() != 2 {
		t.Error("Not Found 2 ")
	}
	conn.Close()
}
