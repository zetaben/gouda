package gouda

type connection interface {
	Open(connectionString string) bool
	Query(r Relation) string
	Close()
}



