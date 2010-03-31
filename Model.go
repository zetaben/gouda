package gouda

import (
	"reflect"
	"fmt"
	"strings"
)
/** Types **/
type ModelInterface interface {
	TableName() string
	Identifier() string
}

type Model struct {
	tablename  string
	identifier string
	attributes map[string]reflect.Type
	object_cache map[int]map[string]Value
	runtype    reflect.Type
	connection *Connection
}

type ModelRelation struct {
	model *Model
	relation *Relation
}

type ModelStore map[string]*Model

var _ModelStore = make(ModelStore)

/** NullModel **/

type NullModel struct{}

func (n NullModel) TableName() string { return "NilTable create a TableName" }

func (n NullModel) Identifier() string { return "Id" }


/** utils **/

func attributes(m interface{}) (map[string]reflect.Type, reflect.Type) {
	var st *reflect.StructType
	var typ reflect.Type
	if _, ok := reflect.Typeof(m).(*reflect.PtrType); ok {
		typ = reflect.Typeof(m).(*reflect.PtrType).Elem()
	} else {
		typ = reflect.Typeof(m)
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

	return ret, typ
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

func (m *Model) Last() interface{} {
	q := NewRelation(m.tablename).Order(strings.ToLower(m.identifier), "desc").First()
	ret := m.connection.Query(q)
	if(ret.Len()<1){return nil}
	v := ret.At(0).(map[string]Value)
	return m.translateObject(v)
}

func (m *Model) First() interface{} {
	q := NewRelation(m.tablename).First()
	ret := m.connection.Query(q)
	if(ret.Len()<1){return nil}
	v := ret.At(0).(map[string]Value)
	return m.translateObject(v)
}

func (m *Model) All() []interface{} {
	q := NewRelation(m.tablename)
	ret := m.connection.Query(q)
	v := make([]interface{}, ret.Len())
	for i := 0; i < ret.Len(); i++ {
		v[i] = m.translateObject(ret.At(i).(map[string]Value))
	}
	return v
}

func (m *Model) Refresh(a interface{}) interface{} {
	st:=reflect.NewValue(a)
	if p,ok:=st.(*reflect.PtrValue);ok {
		st=p.Elem()
	}

	id:=fmt.Sprint(m.getId(st.(*reflect.StructValue)))
	q := NewRelation(m.tablename).Where(m.identifier+" = '"+id+"'").First()
	ret := m.connection.Query(q)
	if(ret.Len()<1){return nil}
	v := ret.At(0).(map[string]Value)
	return m.translateObject(v)
}

func (m *Model) getId(st *reflect.StructValue) int {
	return st.FieldByName(m.identifier).(*reflect.IntValue).Get()
}

func (m *Model) Delete(a interface{}) interface{}{
	st:=reflect.NewValue(a)
	if p,ok:=st.(*reflect.PtrValue);ok {
		st=p.Elem()
	}
	id:=fmt.Sprint(m.getId(st.(*reflect.StructValue)))
	q := NewRelation(m.tablename).Where(m.identifier+" = '"+id+"'").Delete()
	m.connection.Query(q)
	return a
}
func (m *Model) Save(a interface{}) interface{}{
	stv:=reflect.NewValue(a)
	if p,ok:=stv.(*reflect.PtrValue);ok {
		stv=p.Elem()
	}
	st:=stv.(*reflect.StructValue)
	id:=m.getId(st)
	if v,present := m.object_cache[id]; present {
	if up:=m.buildUpdateMap(st,v); len(up) > 0 {
	r:=new(Relation)
	r.Table(m.tablename)
	r.Update(up,m.identifier,id)
	m.connection.Query(r)
	}
	return a
	}

	r:=new(Relation)
	r.Table(m.tablename)
	r.Insert(m.translateMap(st))
	m.connection.Query(r)

	//Ugly Hack to get Last Inserted Id
	q := NewRelation(m.tablename).Order(strings.ToLower(m.identifier), "desc").First()
	ret := m.connection.Query(q)
	if(ret.Len()<1){return nil}
	v := ret.At(0).(map[string]Value)
	m.translateObjectValue(v,st)
	return a
}

func (m *Model) buildUpdateMap(st *reflect.StructValue,old map[string]Value)  map[string]Value {
	ret:=make(map[string]Value)
	for attr,typ := range m.attributes {
		switch typ.(type){
		case *reflect.IntType:
			if tmp:=st.FieldByName(attr).(*reflect.IntValue).Get() ;int(old[strings.ToLower(attr)].Int())!=tmp {
			ret[attr]=SysInt(tmp).Value()
			}
		case *reflect.StringType:
			if tmp:=st.FieldByName(attr).(*reflect.StringValue).Get() ;string(old[strings.ToLower(attr)].String())!=tmp {
			ret[attr]=SysString(tmp).Value()
			}
		}
	}
return ret
}

func (m *Model) translateMap(obj *reflect.StructValue)  map[string]Value {
	ret:=make(map[string]Value)
	for attr,typ := range m.attributes {
		switch typ.(type){
		case *reflect.IntType:
			ret[attr]=SysInt(obj.FieldByName(attr).(*reflect.IntValue).Get()).Value()
		case *reflect.StringType:
			ret[attr]=SysString(obj.FieldByName(attr).(*reflect.StringValue).Get()).Value()
		case nil:
			ret[attr]=new(_Null)
		}
	}
	return ret
}

func (m *Model) translateObject(v map[string]Value) interface{} {
	p := reflect.MakeZero(m.runtype).(*reflect.StructValue)
	return m.translateObjectValue(v,p)
}

func (m *Model) translateObjectValue(v map[string]Value,p *reflect.StructValue) interface{} {
	for lbl, _ := range m.Attributes() {
		vl := v[strings.ToLower(lbl)]
		switch vl.Kind() {
		//TODO MakeZero ??
		case IntKind:
			tmp := reflect.NewValue(1).(*reflect.IntValue)
			tmp.Set(int(vl.Int()))
			p.FieldByName(lbl).SetValue(tmp)
		case StringKind:
			tmp := reflect.NewValue("").(*reflect.StringValue)
			tmp.Set(string(vl.String()))
			p.FieldByName(lbl).SetValue(tmp)
		}
	}
	m.object_cache[int(v[strings.ToLower(m.identifier)].Int())]=v
	return p.Interface()

}

func Refresh(a interface{}) interface{} { return M(a.(ModelInterface)).Refresh(a) }
func Save(a interface{}) interface{} { return M(a.(ModelInterface)).Save(a) }
func Delete(a interface{}) interface{} { return M(a.(ModelInterface)).Delete(a) }

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
	mod.identifier = m.Identifier()
	attr, run := attributes(m)
	mod.attributes = attr
	mod.runtype = run
	mod.connection = conn
	mod.object_cache = make( map[int]map[string]Value)
	(*st)[modelname] = mod
	return mod
}

/** Model RelationLike methods**/

func (m *Model) newRelation() *ModelRelation{
	mr:=new(ModelRelation)
	mr.model=m
	mr.relation=new(Relation)
	mr.relation.Table(m.tablename)
	return mr
}

func (m *Model) Where(x string) *ModelRelation{
	return m.newRelation().Where(x)
}

func (m *Model) Order(x,y string) *ModelRelation{
	return m.newRelation().Order(x,y)
}

func (m *Model) Count(fields ...[]string) int{
	return m.newRelation().Count(fields)
}

/** ModelRelation **/

func (r *ModelRelation) Where(x string) *ModelRelation {
	r.relation.Where(x)
	return r
}

func (r *ModelRelation) Order(x,y string) *ModelRelation {
	r.relation.Order(x,y)
	return r
}

func (r *ModelRelation) First() interface{} {
	q:=r.relation.First()
	ret := r.model.connection.Query(q)
	if(ret.Len()<1){return nil}
	v := ret.At(0).(map[string]Value)
	return r.model.translateObject(v)
}

func (r *ModelRelation) Last() interface{} {
	q:=r.relation.Order(r.model.identifier,"DESC").First()
	ret := r.model.connection.Query(q)
	if(ret.Len()<1){return nil}
	v := ret.At(0).(map[string]Value)
	return r.model.translateObject(v)
}

func (r *ModelRelation) All() []interface{} {
	ret := r.model.connection.Query(r.relation)
	v := make([]interface{}, ret.Len())
	for i := 0; i < ret.Len(); i++ {
		v[i] = r.model.translateObject(ret.At(i).(map[string]Value))
	}
	return v
}


func (r *ModelRelation) Count(fields ...[]string) int {
	q:=r.relation
	if(len(fields)==0){
		field:=make([]string,1)
		field[0]=r.model.identifier
		q=r.relation.Count(field)
	}else{
		q=r.relation.Count(fields[0])
	}

	ret := r.model.connection.Query(q)
	if(ret.Len()<1){return -1}
	v := ret.At(0).(map[string]Value)
	return int(v["_count"].Int())
}
