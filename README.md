# GoUDA : Go Unified Data Accessor

Gouda is an experimental ORM for Go.

## Features (and planned features):
 - Backend agnostic (though at the moment it knows only how to connect to MySQL and XML files ...)
 - Nice syntax for querying (ex : `MyModel.Where(gouda.F("i").Eq("j")).Where(gouda.F("a").Eq("b")).Order("c","ASC").First()`)
 - Not so intrusive...
 - Pure Go

## Installing 

You first need [Thoj Mysql Client Library](http://github.com/thoj/Go-MySQL-Client-Library/)  

	$ git clone git@github.com:zetaben/gouda.git
	$ cd gouda
	$ make
	$ make install

Or use goinstall (UNTESTED YET !)

## Usage 
	
 1. Declare your structs as ModelInterface (by implementing all the interface methods or by embeding gouda.NullModel)
 2. Connect to your database backend (eg : `gouda.OpenMysql("mysql://root:@localhost:3306/test_db")`)
 3. Optional : Register your models using `gouda.GetModelStore().RegisterModel(m)`
 4. Use it !

 Peruse through Model_test.go to view typical usages

## Known limitations 

At the moment Gouda as the following limitations : 

 - Connects only to mysql and XML files
 - limited support for associations (railslike HasMany, BelongsTo)
 - Knows only how to handle int and string attributes !
 - Connections needs to be registered in the store to be usable automatically... 
 - Needs way more documentation ! 

## Acknowlegments

This work is using work or ideas from other go projects : 
  
  - [Thomas Jager's Mysql Client Library](http://github.com/thoj/Go-MySQL-Client-Library/)
  - [Josh Goebel's godatamapper](http://github.com/yyyc514/go_datamapper) for sqlite3
  - [Michael Stephens' gomongo](http://github.com/mikejs/gomongo/blob/master/bson.go) ideas to provide a common Value interface
