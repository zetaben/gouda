package gouda

type Connection interface {
	Open(connectionString string) bool
	Query(r Relation) string
	Close()
}


type ConnectionStore []*Connection

//TODO : Vector Type
var _ConnectionStore=make(ConnectionStore,10)
var i=0


func (cs * ConnectionStore) RegisterConnection(c * Connection) * Connection {
	(*cs)[i]=c
	i++
	return c
}

func (cs * ConnectionStore) Last() *Connection {
	return (*cs)[i]
}

func GetConnectionStore() *ConnectionStore {return &_ConnectionStore }

