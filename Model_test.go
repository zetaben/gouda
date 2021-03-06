package gouda_test

import (
	"gouda"
	"testing"
	"fmt"
	"sort"
	"reflect"
	"os"
)

/*** Test Models ***/
type Personne struct {
	Nom string
	Id  int
	Age int
	gouda.NullModel
}

func (p Personne) TableName() string { return "personne" }

type Car struct {
	Id       int
	Plate    string
	Model    string
	Owner_id int
	gouda.NullModel
}

/*** Test Helping variables ***/

var conn_ok bool = false

/*** Test Functions ***/

func TestTableName(t *testing.T) {
	need_connection()
	var c Car
	Cars := gouda.M(c)
	if Cars.TableName() != "cars" {
		t.Error("Wrong TableName")
	}
}

func TestAttributes(t *testing.T) {
	need_connection()
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
	need_connection()
	p := new(Personne)
	attr := gouda.M(p).AttributesNames()
	names := []string{"Nom", "Id", "Age"}
	sort.StringArray(attr).Sort()
	sort.StringArray(names).Sort()
	if len(names) != len(attr) {
		t.Error("Attributes Names, found  size mismatch : " + fmt.Sprint(len(attr)) + " for " + fmt.Sprint(len(names)))
	}

	if !reflect.DeepEqual(names, attr) {
		t.Error("Can't find Attributes Names, found : " + fmt.Sprint(attr))
	}
	var pp Personne
	attr = gouda.M(pp).AttributesNames()
	sort.StringArray(attr).Sort()
	if len(names) != len(attr) {
		t.Error("Attributes Names, found  size mismatch : " + fmt.Sprint(len(attr)) + " for " + fmt.Sprint(len(names)))
	}

	if !reflect.DeepEqual(names, attr) {
		t.Error("Can't find Attributes Names, found : " + fmt.Sprint(attr))
	}
}


func TestModelName(t *testing.T) {
	need_connection()
	p := new(Personne)
	if fname := gouda.ModelName(p); fname != "Personne-personne" {
		t.Error("wrong name found : " + fname)
	}

	var pp Personne
	if fname := gouda.ModelName(pp); fname != "Personne-personne" {
		t.Error("wrong name found : " + fname)
	}
}

func TestModelFetch(t *testing.T) {
	need_connection()
	p := new(Personne)
	z := gouda.M(p).First().(Personne)
	if z.Id != 1 {
		t.Error("wrong personne found : Id : " + fmt.Sprint(z.Id))
	}

	z = gouda.M(p).Last().(Personne)
	if z.Id != 2 {
		t.Error("wrong personne found : Id : " + fmt.Sprint(z.Id))
	}

	tab := gouda.M(p).All()
	if len(tab) != 2 {
		t.Error("Wrong Fetched Size (" + fmt.Sprint(len(tab)) + ")")
	}
}

func TestModelRelationFetch(t *testing.T) {
	var p Personne
	need_connection()
	Personnes := gouda.M(p)
	toto := Personnes.Where(gouda.F("nom").Eq("toto")).First().(Personne)
	if toto.Nom != "toto" || toto.Id != 1 {
		t.Error("Not Found toto")
	}
	toto = Personnes.Where(gouda.F("nom").Eq("toto")).Last().(Personne)
	if toto.Nom != "toto" || toto.Id != 1 {
		t.Error("Not Found toto")
	}

	totos := Personnes.Where(gouda.F("nom").Eq("toto")).All()

	if len(totos) != 1 {
		t.Fatal("Wrong Fetched Size, fetched :" + fmt.Sprint(len(totos)))
	}
	toto = totos[0].(Personne)
	if toto.Nom != "toto" || toto.Id != 1 {
		t.Error("Not Found toto")
	}
	p = Personnes.Where(gouda.F("id").Gt(1)).First().(Personne)
	if p.Nom != "titi" || p.Id != 2 {
		t.Error("Not Found titi")
	}

	p = Personnes.Where(gouda.F("id").Lt(2)).First().(Personne)
	if p.Nom != "toto" || p.Id != 1 {
		t.Error("Not Found toto")
	}

	p = Personnes.Where(gouda.F("id").GtEq(1)).First().(Personne)
	if p.Nom != "toto" || p.Id != 1 {
		t.Error("Not Found toto")
	}

	p = Personnes.Order("id", "DESC").Where(gouda.F("id").LtEq(2)).First().(Personne)
	if p.Nom != "titi" || p.Id != 2 {
		t.Error("Not Found titi")
	}

	p = Personnes.Where(gouda.F("id").NEq(1)).First().(Personne)
	if p.Nom != "titi" || p.Id != 2 {
		t.Error("Not Found titi")
	}
	totos = Personnes.Where(gouda.F("id").Eq(1).Or(gouda.F("age").IsNull())).All()
	if len(totos) != 2 {
		t.Fatal("Wrong Fetched Size, fetched :" + fmt.Sprint(len(totos)))
	}
}

func TestModelRelationFetchOrder(t *testing.T) {
	var p Personne
	need_connection()
	Personnes := gouda.M(p)
	totos := Personnes.Order("nom", "asc").All()

	if len(totos) != 2 {
		t.Error("Wrong Fetched Size, fetched :" + fmt.Sprint(len(totos)))
	}
	i := 0

	if totos[i].(Personne).Nom != "titi" {
		t.Error("Not Found titi " + fmt.Sprint(totos[i].(Personne)))
	}

	i++

	if totos[i].(Personne).Nom != "toto" {
		t.Error("Not Found toto " + fmt.Sprint(totos[i].(Personne)))
	}
}

func TestModelRelationFetchCount(t *testing.T) {
	var p Personne
	need_connection()
	Personnes := gouda.M(p)
	if Personnes.Count() != 2 {
		t.Error("Counting Personne failed, counted : " + fmt.Sprint(Personnes.Count()))
	}

	if Personnes.Where(gouda.F("nom").Eq("toto")).Count() != 1 {
		t.Error("Counting toto failed, counted : " + fmt.Sprint(Personnes.Count()))
	}

	if Personnes.Count([]string{"age"}) != 1 {
		t.Error("Counting toto age failed, counted : " + fmt.Sprint(Personnes.Count([]string{"age"})))
	}

}

func TestModelRelationRefresh(t *testing.T) {
	var p Personne
	need_connection()
	Personnes := gouda.M(p)
	toto := Personnes.First().(Personne)
	toto = Personnes.Refresh(toto).(Personne)
	if toto.Id != 1 {
		t.Error("Refresh Failed, " + fmt.Sprint(toto))
	}

	toto = Personnes.Refresh(&toto).(Personne)
	if toto.Id != 1 {
		t.Error("Refresh Failed, " + fmt.Sprint(toto))
	}

	toto = gouda.Refresh(&toto).(Personne)
	if toto.Id != 1 {
		t.Error("Refresh Failed, " + fmt.Sprint(toto))
	}
}

func TestModelAssociations(t *testing.T) {
	var p Personne
	var c Car
	need_connection()
	Personnes := gouda.M(p)
	Cars := gouda.M(c)
	Cars.BelongsToKey(Personnes, "owner", "Owner_id")
	Personnes.HasManyKey(Cars, "cars", "Owner_id")
	car := Cars.First().(Car)
	p = Cars.GetAssociated("owner", car).(Personne)
	if p.Nom != "toto" || p.Id != 1 {
		t.Error("Not Found toto")
	}
	cars := Personnes.GetAssociated("cars", p).([]interface{})
	if len(cars) != 2 {
		t.Error("Associated cars not found")
	}
	if cars[0].(Car).Id != 1 {
		t.Error("Associated car not found (first one)")
	}
}

func need_connection() {
	if !conn_ok {
		init_mysql()
		//init_xml()
		conn_ok = true
	}

}


func TestModelInsertOrUpdate(t *testing.T) {
	var p Personne
	need_connection()
	gouda.Save(&Personne{Id: 3, Nom: "test", Age: 12})
	Personnes := gouda.M(p)
	p.Id = 0
	p.Nom = "tata"
	p.Age = 7
	Personnes.Save(&p)
	if p.Id != 4 {
		t.Error("Id not setted, found " + fmt.Sprint(p.Id))
	}
	p.Nom = "plop"
	Personnes.Save(&p)
	if p.Id != 4 {
		t.Error("Save should have updated not inserted!")
	}
	p = gouda.Refresh(&p).(Personne)
	if p.Nom != "plop" {
		t.Error("Not Updated !")
	}
}
func TestModelDelete(t *testing.T) {
	var p Personne
	need_connection()
	Personnes := gouda.M(p)
	Personnes.Delete(Personnes.First())
	if len(Personnes.All()) != 3 {
		t.Error("Not deleted ! counting :" + fmt.Sprint(len(Personnes.All())))
	}
	gouda.Delete(Personnes.First())
	if len(Personnes.All()) != 2 {
		t.Error("Not deleted ! counting :" + fmt.Sprint(len(Personnes.All())))
	}
}

func init_mysql() {
	r, w, err := os.Pipe()
	if err != nil {
		panic("%v", err)
	}

	fmt.Print("Initializing DB... ")
	pid, _ := os.ForkExec("/usr/bin/mysql", []string{"/usr/bin/mysql", "test_db"}, os.Environ(), "/versatile", []*os.File{r, os.Stdout, os.Stderr})
	//	fmt.Fprintln(w,"show tables;");
	fmt.Fprintln(w, "DROP TABLE personne;")
	fmt.Fprintln(w, "DROP TABLE cars;")
	fmt.Fprintln(w, "CREATE TABLE `personne` ( `id` int(11) NOT NULL auto_increment, `nom` varchar(255) default NULL,  age int(3) default NULL,   PRIMARY KEY  (`id`)  );")
	fmt.Fprintln(w, "INSERT INTO `personne` VALUES (1,'toto',23);")
	fmt.Fprintln(w, "INSERT INTO `personne` VALUES (2,'titi',NULL);")
	fmt.Fprintln(w, "CREATE TABLE `cars` ( `id` int(11) NOT NULL auto_increment, `plate` varchar(255) default NULL,   `model` varchar(255) default NULL, owner_id int(11) default NULL,   PRIMARY KEY  (`id`)  );")
	fmt.Fprintln(w, "INSERT INTO `cars` VALUES (1,'123ABC12','Renault',1);")
	fmt.Fprintln(w, "INSERT INTO `cars` VALUES (2,'123CBA12','Traban',1);")
	w.Close()
	os.Wait(pid, 0)
	fmt.Println("Finished!")

	conn := gouda.OpenMysql("mysql://root:@localhost:3306/test_db")
	gouda.GetConnectionStore().RegisterConnection(&conn)
}

func init_xml() {
	conStr := "personnes.xml"
	conn := gouda.OpenXML(conStr).(*gouda.XMLConnector)
	table := conn.CreateTable("personne")
	table.AddAttribute("nom", gouda.StringKind)
	table.AddAttribute("age", gouda.IntKind)
	table.Insert(map[string]gouda.Value{"id": gouda.SysInt(1).Value(), "nom": gouda.SysString("toto").Value(), "age": gouda.SysInt(13).Value()})
	table.Insert(map[string]gouda.Value{"id": gouda.SysInt(2).Value(), "nom": gouda.SysString("titi").Value(), "age": gouda.SysInt(0).Value()})
	table = conn.CreateTable("cars")
	table.AddAttribute("plate", gouda.StringKind)
	table.AddAttribute("model", gouda.StringKind)
	table.AddAttribute("owner_id", gouda.IntKind)
	table.Insert(map[string]gouda.Value{"id": gouda.SysInt(1).Value(), "plate": gouda.SysString("123ABC12").Value(), "model": gouda.SysString("Renault").Value(), "owner_id": gouda.SysInt(1).Value()})
	conn.Close()
	conn2 := gouda.OpenXML(conStr)
	gouda.GetConnectionStore().RegisterConnection(&conn2)

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
