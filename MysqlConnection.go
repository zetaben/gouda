package gouda

import (
	"mysql"
	"strings"
	"fmt"
	"os"
	"strconv"
	"container/vector"
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

	fmt.Println(len(res.ResultSet.Rows))
	var ret vector.Vector
	tmp := make(map[string]Value)
	for rowmap := res.FetchRowMap(); rowmap != nil; rowmap = res.FetchRowMap() {
		fmt.Printf("%#v\n", rowmap)
		fmt.Printf("%#v\n", res.ResultSet.Fields)
		for i := 0; i < len(rowmap); i++ {
//			rowmap[rs.ResultSet.Fields[i].Name] = row.Data[i].Data
			fmt.Println(res.ResultSet.Fields[i].Name)
			var val Value;
			fmt.Println(res.ResultSet.Fields[i].Type)
			switch res.ResultSet.Fields[i].Type {
			case mysql.MYSQL_TYPE_VAR_STRING:
				val=SysString(rowmap[res.ResultSet.Fields[i].Name]).Value()
			case mysql.MYSQL_TYPE_LONG:
				t,_:=strconv.Atoi(rowmap[res.ResultSet.Fields[i].Name])
				val=SysInt(t).Value()
			}
			tmp[res.ResultSet.Fields[i].Name] = val
//			tmp["id"] = val
		}
		ret.Push(tmp)
	}
	fmt.Printf("%#v\n",ret)
	return "plip"
}

func OpenMysql(conStr string) *mysqlConnector {
	db := new(mysqlConnector)
	db.Open(conStr)
	return db
}
