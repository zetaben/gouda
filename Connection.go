package gouda

type Connection interface {
	Open(connectionString string) bool
	Query(r Relation) string
	Close()
}



