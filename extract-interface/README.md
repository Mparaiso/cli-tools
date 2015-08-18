#extract-interface
	
Author:  mparaiso <mparaiso@online.fr>

Year: 2015

license: GPL-3

Extract-interface allows go developers to generate interfaces from structs with methods. 

With Extract-interface , developers no longer need to write interfaces derived from structs by hand.

Extract-interface is written in go

###instal

$ go get github.com/interactiv/cli-tools/extract-interface

###usage

Within a package directory

$ extract-interface

###options

- -dir="": Directory where to find structs to extract interfaces from,
        example: -dir=/My/Directory
- -types="":
        Type Filter , will only extract interfaces from listed types,
        example : -types=MyStuctType1,MyStructType2
