package gouda

import (
	"xml"
	//	"strings"
	"fmt"
	"os"
	"bytes"
	"container/vector"
)

/** Value **/

type XMLConnector struct {
	db     *os.File
	tables map[string]XMLTable
}

type XMLTable map[string]int //Kind

func (e *XMLConnector) Close() {
	e.commit()
	e.db.Close()
}


func (e *XMLConnector) Open(connectionString string) bool {
	f, err := os.Open(connectionString, os.O_RDWR|os.O_CREAT, 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error openning"+fmt.Sprint(err))
		os.Exit(1)
	}
	e.db = f
	e.init()
	return false
}

func xml_attribute(xel xml.StartElement, name string) (xml.Attr, bool) {
	for _, a := range xel.Attr {
		if a.Name.Local == name {
			return a, true
		}
	}
	return xml.Attr{}, false
}

func (e *XMLConnector) init() {
	parser := xml.NewParser(e.db)
	a, err := parser.Token()
	schema := false
	lastType := NullKind
	var lastTable XMLTable
	for ; err == nil; a, err = parser.Token() {
		switch a.(type) {
		case xml.StartElement:
			xel := a.(xml.StartElement)
			el := xel.Name.Local
			//fmt.Println(el)
			if el == "schema" {
				schema = true
			}
			if schema {
				switch el {
				case "table":
					if attr, found := xml_attribute(xel, "name"); found {
					//	fmt.Println("Found Table ! " + attr.Value)
						lastTable = e.CreateTableWithoutId(attr.Value)
					} else {
					//	fmt.Println(os.Stderr, "Missing required attribute name on <table>")
					}
				case "attribute":
					if attr, found := xml_attribute(xel, "type"); found {
					//	fmt.Println("Found Attribute ! on" + fmt.Sprint(lastTable) + " " + attr.Value)
						//fmt.Println(xel)
						lastType = typeFromName(attr.Value)
					} else {
						fmt.Println(os.Stderr, "Missing required attribute type on <attribute>")
					}

				}
			}
		case xml.CharData:
			if schema && lastType != NullKind {
				b := bytes.NewBuffer(a.(xml.CharData))
				lastTable.AddAttribute(b.String(), lastType)
				lastType = NullKind
			}
		case xml.EndElement:
			el := a.(xml.EndElement).Name.Local
			//fmt.Println("/" + el)
			if el == "schema" {
				schema = false
			}
		}
	}
	if _,oser:=err.(os.Error); err != nil && !oser  {
		fmt.Printf("%T ", err)
		fmt.Println(err)
	}
}

func (e *XMLConnector) commit() {
	e.db.Seek(0, 0)
	e.db.Truncate(0)
	fmt.Fprintln(e.db, "<?xml version=\"1.0\" ?>")
	fmt.Fprintln(e.db, "<db>")
	fmt.Fprintln(e.db, "<schema>")
	for table, tabledesc := range e.tables {
//		fmt.Println("Comitting table " + table + " struct " + fmt.Sprint(tabledesc))
		fmt.Fprintf(e.db, "<table name=\"%s\">\n", table)
		for key, typ := range tabledesc {
			fmt.Fprintf(e.db, "<attribute type=\"%s\">%s</attribute>\n", typeName(typ), key)
		}
		fmt.Fprintln(e.db, "</table>")
	}
	fmt.Fprintln(e.db, "</schema>")
	fmt.Fprint(e.db, "</db>")

}

//TODO type ValueKind ??
func typeName(kind int) string {
	switch kind {
	case IntKind:
		return "int"
	case StringKind:
		return "string"
	}
	return "UNKNOWN"
}

func typeFromName(kind string) int {
	switch kind {
	case "int":
		return IntKind
	case "string":
		return StringKind
	}
	return NullKind
}


func (e *XMLConnector) CreateTableWithoutId(table string) XMLTable {
	e.tables[table] = make(XMLTable)
	return e.tables[table]
}

func (e *XMLConnector) CreateTable(table string) XMLTable {
	a := e.CreateTableWithoutId(table)
	a.AddAttribute("id", IntKind)
	return a
}

func (e *XMLConnector) Table(table string) XMLTable {
	tab, ok := e.tables[table]
	if ok {
		return tab
	}
	return nil
}

func (e *XMLTable) AddAttribute(name string, typ int) {
	(*e)[name] = typ
}
func (e XMLTable) Attributes() map[string]int { return e }

func (e *XMLConnector) Query(r *Relation) *vector.Vector {
	return new(vector.Vector)
}

func OpenXML(conStr string) Connection {
	db := (new(XMLConnector))
	db.tables = make(map[string]XMLTable)
	db.Open(conStr)
	return db
}
