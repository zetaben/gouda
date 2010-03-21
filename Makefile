#*****************************
#This creates my package 'mysqlgo'

include $(GOROOT)/src/Make.$(GOARCH)

TARG=gouda
GOFILES=\
Relation.go\
Connection.go\
MysqlConnection.go\
Value.go\
Model.go\

include $(GOROOT)/src/Make.pkg
