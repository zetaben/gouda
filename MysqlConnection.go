package gouda

import (
	"mysql"
	"strings"
	"fmt"
)

type mysqlConnector struct {
	conn * mysql.MySQLInstance
}


func (*mysqlConnector ) Open(conStr string) bool {
	fmt.Println("Plop")
	tab:=strings.Split(conStr,"/",0)
	fmt.Println(tab[2])
	return false
}

func (*mysqlConnector ) Query(conStr string) string {
	fmt.Println("Plip")
	return "plip"
}

func OpenMysql(conStr string) *mysqlConnector {
	db:= new(mysqlConnector)
	db.Open(conStr)
	return db
}
