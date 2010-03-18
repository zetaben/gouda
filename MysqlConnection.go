package gouda

import (
	"mysql"
	"strings"
	"fmt"
	"os"
)

type mysqlConnector struct {
	conn *mysql.MySQLInstance
}


func (e *mysqlConnector) Open(conStr string) bool {
	fmt.Println("Plop")
	tab := strings.Split(conStr, "/", 0)
	db := tab[len(tab)-1]
	tab2 := strings.Split(tab[2], "@", 2)
	tab = strings.Split(tab2[0], ":", 0)
	fmt.Println(tab)
	fmt.Println(tab2)
	user := tab[0]
	pass := tab[1]

	dbh, err := mysql.Connect("tcp", "", tab2[1], user, pass, "")
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	e.conn = dbh
	e.conn.Use(db)
	return false
}

func (e *mysqlConnector) Query(r *Relation) string {
	res, err := e.conn.Query(r.Sql())
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	fmt.Println(res)
	fmt.Println(res.FieldCount)

	fmt.Println(res)

	for rowmap := res.FetchRowMap(); rowmap != nil; rowmap = res.FetchRowMap() {
		fmt.Printf("%#v\n", rowmap)
		fmt.Printf("%#v\n", res.ResultSet.Fields)

	}

	return "plip"
}

func OpenMysql(conStr string) *mysqlConnector {
	db := new(mysqlConnector)
	db.Open(conStr)
	return db
}
