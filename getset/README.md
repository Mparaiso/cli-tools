#getset
	
Author:  mparaiso <mparaiso@online.fr>

Year: 2015

license: GPL-3

GetSet generates getters and setters for your structs. With getset , go developers
no longer need to write their getters and setters by hand. getset is written in go

###installation : 

go get github.com/interactiv/cli-tools/getset

###usage:

within a package directory

$ getset

####options:

- -dir="": package directory.defaults to working directory
- -file=false: write to file, true of false, false by default, 
	if false then will write to the standard out stdout.
- -type="": type from which getters and setters will be generated, declare 
	multiple types by separating them with a comma, for instance: -type=Foo,Bar,Baz