#*****************************
#This creates my package 'mysqlgo'

include $(GOROOT)/src/Make.$(GOARCH)

TARG=gouda
GOFILES=\
Relation.go\

include $(GOROOT)/src/Make.pkg
