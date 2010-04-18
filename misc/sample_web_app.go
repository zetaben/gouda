package main

/** 
* Gouda Sample App
* Using Web.go 
* Constructed using web.go sample app
* 
SQL Structure : 

CREATE TABLE `cars` (
  `id` int(11) NOT NULL auto_increment,
  `plate` varchar(255) default NULL,
  `model` varchar(255) default NULL,
  `owner_id` int(11) default NULL,
  PRIMARY KEY  (`id`)
) DEFAULT CHARSET=utf8;

CREATE TABLE `personne` (
  `id` int(11) NOT NULL auto_increment,
  `nom` varchar(255) default NULL,
  `age` int(3) default NULL,
  PRIMARY KEY  (`id`)
) DEFAULT CHARSET=utf8;

* compile & run using : 
* 8g sample_web_app.go && 8l  -o sample_web_app sample_web_app.8  && ./sample_web_app
*/

import (
	"web"
	"gouda"
	"fmt"
	"strconv"
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

func person_detail(val string) string {
	var p Personne
	Personnes := gouda.M(p)
	//Cars:=gouda.M(c)
	i, _ := strconv.Atoi(val)
	p = Personnes.Where(gouda.F("id").Eq(i)).First().(Personne)
	ret:="<div>"
	ret+="<p> Nom : "+p.Nom+"</p>"
	ret+="<p> Age : "+fmt.Sprint(p.Age)+"</p>"
	ret+="</div>"
	cars:=Personnes.GetAssociated("cars",p).([]interface{})
	ret+="<ul>"
	for _,c:=range cars {
	ret+="<li>"+c.(Car).Plate+"</li>"
	}
	ret+="</ul>"
	ret+="<a href=\"/edit_person/"+fmt.Sprint(p.Id)+"\"/>Editer</a>"
	return ret
}

func person_edit(val string) string {
	var p Personne
	Personnes := gouda.M(p)
	i, _ := strconv.Atoi(val)
	p = Personnes.Where(gouda.F("id").Eq(i)).First().(Personne)
	ret:="<div><form action=\"/update_person\" method=\"post\">"
	ret+="<p> Nom : <input type=\"text\" name=\"nom\" value=\""+p.Nom+"\"></p>"
	ret+="<p> Age : <input type=\"text\" name=\"age\" value=\""+fmt.Sprint(p.Age)+"\"></p>"
	ret+="<input type=\"hidden\" name=\"id\" value=\""+fmt.Sprint(p.Id)+"\"><input type=\"submit\" /></form></div>"
	return ret
}

func person_delete(ctx *web.Context,val string) {
	var p Personne
	var c Car
	Personnes := gouda.M(p)
	Cars:=gouda.M(c)
	i, _ := strconv.Atoi(val)
	p = Personnes.Where(gouda.F("id").Eq(i)).First().(Personne)
	cars:=Personnes.GetAssociated("cars",p).([]interface{})
	for _,c:=range cars {
	Cars.Delete(c.(Car))
	}
	Personnes.Delete(p)
	ctx.Redirect(302, "/")
}
func update_person(ctx *web.Context) {
	var p Personne
	Personnes := gouda.M(p)
	i, _ := strconv.Atoi(ctx.Request.Params["id"][0])
	p = Personnes.Where(gouda.F("id").Eq(i)).First().(Personne)
	p.Nom = ctx.Request.Params["nom"][0]
	p.Age, _ = strconv.Atoi(ctx.Request.Params["age"][0])
	Personnes.Save(p)
	ctx.Redirect(302, "/person/"+fmt.Sprint(p.Id))
}


func persons() string {
	var p Personne
	Personnes := gouda.M(p)
	coll := Personnes.Order("id", "ASC").All()
	ret := "<h1>Personnes</h1>"
	ret += "<table>"
	ret += "<tr><th>Id</th><th>Nom</th><th>Age</th><th>Action</th></tr>"
	for _, r := range coll {
		p = r.(Personne)
		ret += "<tr><td><a href=\"/person/" + fmt.Sprint(p.Id) + "\">" + fmt.Sprint(p.Id) + "</a></td><td>" + p.Nom + "</td><td>" + fmt.Sprint(p.Age) + "</td><td><a href=\"/delete_person/" + fmt.Sprint(p.Id) + "\">suppr</a></td></tr>"
	}
	ret += "</table><a href=\"/new_person/\">Ajouter Person</a>"
	return ret
}

func hello(val string) string { return "hello " + val }

func need_connection() (gouda.Connection) {
	conn := gouda.OpenMysql("mysql://root:@localhost:3306/test_db")
	gouda.GetConnectionStore().RegisterConnection(&conn)
	return conn
}

func new_person() string {
	return "<form method=\"post\" action=\"/person\">Nom : <input type=\"text\" name =\"nom\"/><br />Age : <input type=\"text\" name =\"age\"/><input type=\"submit\"></form>"
}

func create_person(ctx *web.Context) {
	var p Personne
	Personnes := gouda.M(p)
	p.Id = 0
	p.Nom = ctx.Request.Params["nom"][0]
	p.Age, _ = strconv.Atoi(ctx.Request.Params["age"][0])
	Personnes.Save(p)
	ctx.Redirect(302, "/")
}

func main() {

	conn:=need_connection()
	var p Personne
	var c Car
	Personnes := gouda.M(p)
	Cars := gouda.M(c)
	Cars.BelongsToKey(Personnes, "owner", "Owner_id")
	Personnes.HasManyKey(Cars, "cars", "Owner_id")

	//web.Get("/(.*)", hello)
	web.Get("/new_person/", new_person)
	web.Get("/persons/", persons)
	web.Get("/", persons)
	web.Get("/person/(.*)", person_detail)
	web.Get("/edit_person/(.*)", person_edit)
	web.Get("/delete_person/(.*)", person_delete)
	web.Post("/person", create_person)
	web.Post("/update_person", update_person)
	web.Run("0.0.0.0:9999")
	conn.Close()
}
