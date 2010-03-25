package gouda

import (
	"reflect"
	"fmt"
	"strings"
)
/** Types **/
type ModelInterface interface {
	TableName() string
}

type Model struct {
	tablename  string
	attributes map[string]reflect.Type
	runtype  reflect.Type
	connection *Connection
}

type ModelStore map[string]*Model

var _ModelStore = make(ModelStore)

/** utils **/

func attributes(m interface{}) (map[string]reflect.Type,reflect.Type) {
	var st *reflect.StructType
	var typ reflect.Type
	if _, ok := reflect.Typeof(m).(*reflect.PtrType); ok {
		typ = reflect.Typeof(m).(*reflect.PtrType).Elem()
	} else {
		typ =reflect.Typeof(m);
	}
		st = typ.(*reflect.StructType)

	//fmt.Println(st.NumField())

	ret := make(map[string]reflect.Type)

	for i := 0; i < st.NumField(); i++ {
		p := st.Field(i)
		//fmt.Println(p.Name)
		if !p.Anonymous {
			ret[p.Name] = p.Type
		}
	}

	return ret,typ
}

/** Model **/
func (m Model) TableName() string { return m.tablename }


func (m *Model) Attributes() map[string]reflect.Type {
	return m.attributes
}

func (m *Model) AttributesNames() (ret []string) {
	ret = make([]string, len(m.attributes))
	i := 0
	for k, _ := range m.attributes {
		ret[i] = k
		i++
	}
	return ret
}

func (m *Model) First() interface{} {
	q := NewRelation(m.tablename).First()
	ret := m.connection.Query(q)
	v := ret.At(0).(map[string]Value)
	fmt.Println(ret)
	p:=reflect.MakeZero(m.runtype).(*reflect.StructValue)
	for lbl, _ := range m.Attributes() {
//		fmt.Println(lbl)
//		fmt.Println(typ)
		vl := v[strings.ToLower(lbl)]
		switch vl.Kind() {
		case IntKind:
			tmp:=reflect.NewValue(1).(*reflect.IntValue)
			tmp.Set(int(vl.Int()))
//			fmt.Println("FIL "+lbl)
//			fmt.Println(reflect.MakeZero(m.runtype))
			p.FieldByName(lbl).SetValue(tmp)
//			fmt.Println(vl.Int());
		case StringKind:
			fmt.Println(vl.String());
			tmp:=reflect.NewValue("").(*reflect.StringValue)
			tmp.Set(string(vl.String()))
//			fmt.Println("FIL "+lbl)
//			fmt.Println(reflect.MakeZero(m.runtype))
			p.FieldByName(lbl).SetValue(tmp)
		}
	}
//	fmt.Printf("%#v\n",p)
	return p.Interface()
}


/** ModelInterface **/

func ModelName(m ModelInterface) (ret string) {
	t := reflect.Typeof(m).String()
	tab := strings.Split(t, ".", 0)
	return tab[len(tab)-1] + "-" + m.TableName()
}

func M(m ModelInterface) *Model {
	modelname := ModelName(m)
	if model, present := _ModelStore[modelname]; present {
		return model
	}
	return GetModelStore().RegisterModel(m)

}


/** ModelStore **/

func GetModelStore() *ModelStore { return &_ModelStore }

func (st *ModelStore) RegisterModel(m ModelInterface) *Model {
	return st.RegisterModelWithConnection(m, GetConnectionStore().Last())
}
func (st *ModelStore) RegisterModelWithConnection(m ModelInterface, conn *Connection) *Model {
	modelname := ModelName(m)
	mod := new(Model)
	mod.tablename = m.TableName()
	attr,run :=attributes(m)
	mod.attributes =  attr
	mod.runtype =  run
	mod.connection = conn
	(*st)[modelname] = mod
	return mod
}
