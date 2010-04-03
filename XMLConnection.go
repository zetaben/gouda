package gouda

import (
	"xml"
	"strconv"
	"fmt"
	"os"
	"bytes"
	"sort"
	"math"
	"strings"
	"container/vector"
)

type XMLConnector struct {
	db     *os.File
	tables map[string]*XMLTable
}

type ValueVector struct {
	sort      vector.StringVector
	order_asc vector.StringVector
	vector.Vector
}

func (v *ValueVector) Less(i, j int) bool {
	ii := v.At(i).(map[string]Value)
	jj := v.At(j).(map[string]Value)
	for i := 0; i < v.sort.Len(); i++ {
		sort := strings.ToLower(v.sort.At(i))
		order_asc := strings.ToUpper(v.order_asc.At(i)) == "ASC"
		switch ii[sort].Kind() {
		case IntKind:
			if ii[sort].Int() == jj[sort].Int() {
				continue
			}
		case StringKind:
			if ii[sort].String() == jj[sort].String() {
				continue
			}
		}

		if order_asc {
			switch ii[sort].Kind() {
			case IntKind:
				return ii[sort].Int() < jj[sort].Int()
			case StringKind:
				return ii[sort].String() < jj[sort].String()
			}
		} else {
			switch ii[sort].Kind() {
			case IntKind:
				return ii[sort].Int() > jj[sort].Int()
			case StringKind:
				return ii[sort].String() > jj[sort].String()
			}

		}
	}
	return false
}

type XMLTable struct {
	schema map[string]int
	data   *ValueVector
}

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
	var lastAttr string = ""
	var lastTable *XMLTable
	var lastItem map[string]Value
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
			} else {
				switch el {
				case "tabledata":
					if attr, found := xml_attribute(xel, "name"); found {
						lastTable = e.Table(attr.Value)
					} else {
						fmt.Println(os.Stderr, "Missing required attribute name on <tabledata>")
					}
				case "item":
					lastItem = make(map[string]Value)

				case "value":
					if attr, found := xml_attribute(xel, "name"); found {
						lastAttr = attr.Value
					} else {
						fmt.Println(os.Stderr, "Missing required attribute name on <tabledata>")
					}
				}
			}
		case xml.CharData:
			if schema && lastType != NullKind {
				b := bytes.NewBuffer(a.(xml.CharData))
				lastTable.AddAttribute(b.String(), lastType)
				lastType = NullKind
			} else if lastAttr != "" {
				b := bytes.NewBuffer(a.(xml.CharData))
				switch lastTable.Attributes()[lastAttr] {
				case StringKind:
					lastItem[lastAttr] = SysString(b.String()).Value()
				case IntKind:
					i, _ := strconv.Atoi(b.String())
					lastItem[lastAttr] = SysInt(i).Value()

				}
				lastAttr = ""

			}
		case xml.EndElement:
			el := a.(xml.EndElement).Name.Local
			//fmt.Println("/" + el)
			if el == "schema" {
				schema = false
			}

			if schema {
				if el == "table" {
					lastTable = nil
				}
			} else {
				switch el {
				case "tabledata":
					lastTable = nil
				case "item":
					lastTable.data.Push(lastItem)
				}
			}
		}
	}
	if _, oser := err.(os.Error); err != nil && !oser {
		fmt.Printf("%T ", err)
		fmt.Println(err)
	}
}

func (e *XMLConnector) Commit() { e.commit() }
func (e *XMLConnector) commit() {
	e.db.Seek(0, 0)
	e.db.Truncate(0)
	fmt.Fprintln(e.db, "<?xml version=\"1.0\" ?>")
	fmt.Fprintln(e.db, "<db>")
	fmt.Fprintln(e.db, "<schema>")
	for table, tabledesc := range e.tables {
		//		fmt.Println("Comitting table " + table + " struct " + fmt.Sprint(tabledesc))
		fmt.Fprintf(e.db, "<table name=\"%s\">\n", table)
		for key, typ := range tabledesc.Attributes() {
			fmt.Fprintf(e.db, "<attribute type=\"%s\">%s</attribute>\n", typeName(typ), key)
		}
		fmt.Fprintln(e.db, "</table>")
	}
	fmt.Fprintln(e.db, "</schema>")
	fmt.Fprintln(e.db, "<data>")
	for table, tabledesc := range e.tables {
		fmt.Fprintf(e.db, "<tabledata name=\"%s\">\n", table)
		for i := 0; i < tabledesc.Data().Len(); i++ {
			values := tabledesc.Data().At(i).(map[string]Value)
			fmt.Fprintln(e.db, "<item>")

			for key, typ := range tabledesc.Attributes() {
				fmt.Fprintf(e.db, "<value name=\"%s\">", key)
				switch typ {
				case IntKind:
					fmt.Fprintf(e.db, "%d", int(values[key].Int()))
				case StringKind:
					xml.Escape(e.db, bytes.NewBufferString(string(values[key].String())).Bytes())
				}
				fmt.Fprintf(e.db, "</value>\n")
			}

			fmt.Fprintln(e.db, "</item>")
		}
		fmt.Fprintln(e.db, "</tabledata>")
	}
	fmt.Fprintln(e.db, "</data>")
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


func (e *XMLConnector) CreateTableWithoutId(table string) *XMLTable {
	tmp := new(XMLTable)
	e.tables[table] = tmp
	tmp.schema = make(map[string]int)
	tmp.data = new(ValueVector)
	return tmp
}

func (e *XMLConnector) CreateTable(table string) *XMLTable {
	a := e.CreateTableWithoutId(table)
	a.AddAttribute("id", IntKind)
	return a
}

func (e *XMLConnector) Table(table string) *XMLTable {
	tab, ok := e.tables[table]
	if ok {
		return tab
	}
	return nil
}

func (e *XMLTable) AddAttribute(name string, typ int) {
	e.schema[name] = typ
}
func (e *XMLTable) Attributes() map[string]int {
	return e.schema
}
func (e *XMLTable) Data() *ValueVector { return e.data }

//unprotected
func (e *XMLTable) Insert(val map[string]Value) {
	e.data.Push(val)
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func copyVect(v *ValueVector) *ValueVector {
	ret := new(ValueVector)

	for i := 0; i < v.Len(); i++ {
		ret.Push(v.At(i))
	}
	return ret
}

func (e *XMLConnector) match(conds vector.Vector, row map[string]Value) bool {

	if conds.Len() == 0 {
		return true
	} //fast out

	for i := range conds {
		cond := conds.At(i).(*Condition)

		switch cond.operand {
		case NULL:
			return false
		case EQUAL:
			if !Equal(row[cond.field], cond.value) {
				return false
			}
		case ISNOTNULL:
			if isnull(row[cond.field]) {
				return false
			}
		case ISNULL:
			if !isnull(row[cond.field]) {
				return false
			}

		default:
			fmt.Fprintln(os.Stderr, "Unknown cond "+fmt.Sprint(cond))
			return false
		}
	}
	return true
}

func isnull(v Value) bool {
	switch v.Kind() {
	case IntKind:
		return int(v.Int()) == 0
	case StringKind:
		return string(v.String()) == ""
	}
	return false
}


func (e *XMLConnector) Query(r *Relation) *vector.Vector {
	ret := new(vector.Vector)

	//	fmt.Println(r)
	switch r.kind {
	case SELECT:
		dat := copyVect(e.tables[r.table].data)
		if r.order_field.Len() > 0 {
			dat.sort = r.order_field
			dat.order_asc = r.order_direction
			sort.Sort(dat)
			//			fmt.Println(dat)
		}

		limit := math.MaxInt32
		if r.limit_count > 0 {
			limit = r.limit_offset + r.limit_count
		}
		//		fmt.Println(r.limit_offset)
		//		fmt.Println(limit)
		found := 0
		for i := 0; i < dat.Len(); i++ {
			tmp := dat.At(i).(map[string]Value)
			if e.match(r.conditions, tmp) {
				if found >= r.limit_offset {
					ret.Push(tmp)
				}
				found++
				if found == limit {
					return ret
				}
			}
		}
	case COUNT:
		for i := range r.attributes {
			k := r.attributes.At(i)
			r.Where(F(k).IsNotNull())
		}
		dat := e.tables[r.table].data
		count := 0
		for i := 0; i < dat.Len(); i++ {
			tmp := dat.At(i).(map[string]Value)
			if e.match(r.conditions, tmp) {
				count++
			}
		}
		re := make(map[string]Value)
		re["_count"] = SysInt(count).Value()
		ret.Push(re)
	case INSERT:
		ins := make(map[string]Value)
		for k, v := range r.values {
			ins[strings.ToLower(k)] = v
		}
		if _, ok := ins["id"]; ok {
			ins["id"] = SysInt(e.tables[r.table].data.Len() + 1).Value()
		}
		e.tables[r.table].data.Push(ins)
	case DELETE:
		var list vector.IntVector
		dat := e.tables[r.table].data
		limit := math.MaxInt32
		if r.limit_count > 0 {
			limit = r.limit_offset + r.limit_count
		}
		found := 0
		for i := 0; i < dat.Len(); i++ {
			tmp := dat.At(i).(map[string]Value)
			if e.match(r.conditions, tmp) {
				if found >= r.limit_offset {
					list.Push(i)
				}
				found++
				if found == limit {
					return ret
				}
			}
		}
		for i := range list {
			e.tables[r.table].data.Delete(list.At(i))
		}

	case UPDATE:
		dat := e.tables[r.table].data
		limit := math.MaxInt32
		if r.limit_count > 0 {
			limit = r.limit_offset + r.limit_count
		}
		found := 0
		for i := 0; i < dat.Len(); i++ {
			tmp := dat.At(i).(map[string]Value)
			if e.match(r.conditions, tmp) {
				if found >= r.limit_offset {
					for k, v := range r.values {
						tmp[strings.ToLower(k)] = v
					}
				}
				found++
				if found == limit {
					return ret
				}
			}
		}
	}
	return ret
}

func OpenXML(conStr string) Connection {
	db := (new(XMLConnector))
	db.tables = make(map[string]*XMLTable)
	db.Open(conStr)
	return db
}
