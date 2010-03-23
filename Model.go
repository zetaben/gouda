package gouda

import (
	"reflect"
//	"fmt"
	"strings"
)
/** Types **/
type ModelInterface interface {
	TableName() string
}

type Model struct {
	tablename  string
	attributes map[string]reflect.Type
	connection *Connection
}

type ModelStore map[string]*Model

var _ModelStore = make(ModelStore)

/** utils **/

func attributes(m interface{}) map[string]reflect.Type {
	var st *reflect.StructType
	if _, ok := reflect.Typeof(m).(*reflect.PtrType); ok {
		st = reflect.Typeof(m).(*reflect.PtrType).Elem().(*reflect.StructType)
	} else {
		st = reflect.Typeof(m).(*reflect.StructType)
	}

	//fmt.Println(st.NumField())

	ret := make(map[string]reflect.Type)

	for i := 0; i < st.NumField(); i++ {
		p := st.Field(i)
		//fmt.Println(p.Name)
		if !p.Anonymous {
			ret[p.Name] = p.Type
		}
	}

	return ret
}

/** Model **/
func (m Model) TableName() string { return m.tablename }



func (m *Model) Attributes() map[string]reflect.Type {
	return m.attributes
}

func (m *Model) AttributesNames() (ret  []string) {
	ret=make([]string,len(m.attributes))
	i:=0
	for k,_:=range m.attributes {
	ret[i]=k
	i++
	}
	return ret
}

func ModelName(m ModelInterface ) (ret string) {
	t:=reflect.Typeof(m).String()
	tab:=strings.Split(t,".",0)
	return tab[len(tab)-1]+"-"+m.TableName()
}

/** ModelInterface **/

func M(m ModelInterface) *Model {
	modelname:=ModelName(m)
	if model,present:=_ModelStore[modelname]; present {
	return model
	}
/*	mod := new(Model)
	mod.tablename = m.TableName()
	mod.attributes = attributes(m)
	_ModelStore[modelname]=mod
	*/
	return GetModelStore().RegisterModel(m)

}


/** ModelStore **/

func GetModelStore() *ModelStore {
	return &_ModelStore;
}

func (st *ModelStore) RegisterModel(m ModelInterface) *Model {
return st.RegisterModelWithConnection(m,GetConnectionStore().Last())
}
func (st *ModelStore) RegisterModelWithConnection(m ModelInterface,conn * Connection) *Model {
	modelname:=ModelName(m)
	mod := new(Model)
	mod.tablename = m.TableName()
	mod.attributes = attributes(m)
	mod.connection = conn
	(*st)[modelname]=mod
	return mod
}

