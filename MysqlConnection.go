package gouda

import (
	"mysql"
	"strings"
	"fmt"
	"os"
	"strconv"
	"container/vector"
)

/** Value **/

func mysql_string(v Value) string {
	switch v.(type) {
	case *_Int:
		return fmt.Sprint(int(v.Int()))
	case *_String:
		return string(v.String())
	}
	return "UNRECHEABLE"
}

func (c *Condition) mysql_string() string {

	ret := " " + c.field + " "

	switch c.operand {

	case EQUAL:
		ret += " = '" + mysql_string(c.value) + "' "
	case NOTEQUAL:
		ret += " != '" + mysql_string(c.value) + "' "
	case GREATER:
		ret += " > '" + mysql_string(c.value) + "' "
	case LOWER:
		ret += " < '" + mysql_string(c.value) + "' "
	case GREATEROREQUAL:
		ret += " >= '" + mysql_string(c.value) + "' "
	case LOWEROREQUAL:
		ret += " <= '" + mysql_string(c.value) + "' "
	case ISNULL:
		ret += " IS NULL "
	case ISNOTNULL:
		ret += " IS NOT NULL "

	default:
		ret += "UNKWON OPERAND " + c.String()

	}

	return ret
}


type MysqlConnector struct {
	conn *mysql.MySQLInstance
}

func (e *MysqlConnector) Close() { e.conn.Quit() }


func (e *MysqlConnector) Open(connectionString string) bool {
	tab := strings.Split(connectionString, "/", 0)
	db := tab[len(tab)-1]
	tab2 := strings.Split(tab[2], "@", 2)
	tab = strings.Split(tab2[0], ":", 0)
	//	fmt.Println(tab)
	//	fmt.Println(tab2)
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


func (e *MysqlConnector) Query(r *Relation) *vector.Vector {
	query := mysql_query(r)
	res, err := e.conn.Query(query)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	ret := new(vector.Vector)
	if r.kind != SELECT && r.kind != COUNT {
		return ret
	}
	//	fmt.Println(res)
	//	fmt.Println(res.FieldCount)

	//	fmt.Println(len(res.ResultSet.Rows))
	for rowmap := res.FetchRowMap(); rowmap != nil; rowmap = res.FetchRowMap() {
		tmp := make(map[string]Value)
		//		fmt.Printf("%#v\n", rowmap)
		//		fmt.Printf("%#v\n", res.ResultSet.Fields)
		for i := 0; i < len(rowmap); i++ {
			//			rowmap[rs.ResultSet.Fields[i].Name] = row.Data[i].Data
			//			fmt.Println(res.ResultSet.Fields[i].Name)
			var val Value
			//			fmt.Println(res.ResultSet.Fields[i].Type)
			switch res.ResultSet.Fields[i].Type {
			case mysql.MYSQL_TYPE_VAR_STRING:
				val = SysString(rowmap[res.ResultSet.Fields[i].Name]).Value()
			case mysql.MYSQL_TYPE_LONG, mysql.MYSQL_TYPE_LONGLONG:
				t, _ := strconv.Atoi(rowmap[res.ResultSet.Fields[i].Name])
				val = SysInt(t).Value()
			}
			tmp[res.ResultSet.Fields[i].Name] = val
			//			tmp["id"] = val
		}
		ret.Push(tmp)
	}
	//	fmt.Printf("%#v\n",ret)
	return ret
}

func OpenMysql(conStr string) Connection {
	db := (new(MysqlConnector))
	db.Open(conStr)
	return db
}


func mysql_query(r *Relation) (sql string) {

	switch r.kind {
	case INSERT:
		sql = "INSERT INTO " + r.table + " ("
		i := 0
		values := "( "
		for k, v := range r.values {
			i++
			sql += strings.ToLower(k)
			values += "'" + mysql_string(v) + "'"
			if i < len(r.values) {
				sql += ", "
				values += ", "
			}
		}

		sql += " ) VALUES " + values + " ) "

	case UPDATE:

		sql = "UPDATE " + r.table + " "
		sql += "SET "
		i := 0
		for k, v := range r.values {
			i++
			sql += strings.ToLower(k)
			sql += " = "
			sql += "'" + mysql_string(v) + "'"
			if i < len(r.values) {
				sql += ", "
			}
		}

		sql += " WHERE " + strings.ToLower(r.identifier_field) + " = '" + fmt.Sprint(r.id) + "'"
	case DELETE:
		sql = "DELETE FROM " + r.table
		if r.conditions.Len() > 0 {
			sql += " WHERE ( "
			for _, ss := range r.conditions {
				sql += ss.(*Condition).mysql_string()
				if ss != r.conditions.Last() {
					sql += " ) AND ( "
				}
			}
			sql += " )"
		}
	default:

		sql = "Select "
		if r.kind == COUNT && r.attributes.Len() > 0 {
			sql += " COUNT( "
			for _, ss := range r.attributes {
				sql += ss
				if ss != r.attributes.Last() {
					sql += ", "
				}
			}
			sql += " ) as _count"
		} else {
			sql += " * "
		}

		sql += " from " + r.table
		if r.conditions.Len() > 0 {
			sql += " where ( "
			for _, ss := range r.conditions {
				sql += ss.(*Condition).mysql_string()
				if ss != r.conditions.Last() {
					sql += " ) AND ( "
				}
			}
			sql += " )"
		}
		if r.order_field.Len() > 0 {
			sql += " ORDER BY "
			for i, ss := range r.order_field {
				sql += ss + " " + r.order_direction[i]
			}
		}

		if r.limit_count > 0 {
			sql += " LIMIT " + fmt.Sprint(r.limit_offset) + ", " + fmt.Sprint(r.limit_count)
		}
	}
	sql += ";"
	//	fmt.Println(sql)
	return
}
