package gouda

import (
	"container/vector"
//	"os"
//	"fmt"
)

type Connection interface {
	Open(connectionString string) bool
	Query(r *Relation) *vector.Vector
	Close()
}
//TODO : *Connection or Connection ? 

type ConnectionStore []*Connection

//TODO : Vector Type
var _ConnectionStore=make(ConnectionStore,10)
var i=0


func (cs * ConnectionStore) RegisterConnection(c *Connection) * Connection {
	(*cs)[i]=c
	i++
	return c
}

func (cs * ConnectionStore) Last() *Connection {
	if i==0 { panic("No Connection defined !")}
	return (*cs)[i-1]
}

func GetConnectionStore() *ConnectionStore {return &_ConnectionStore }

