//    extract-interface version 0.1
//    Copyright (C) 2015  mparaiso <mparaiso@online.fr>
//
//    This program is free software: you can redistribute it and/or modify
//    it under the terms of the GNU General Public License as published by
//    the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.

//    This program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU General Public License for more details.

//    You should have received a copy of the GNU General Public License
//    along with this program.  If not, see <http://www.gnu.org/licenses/>

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"runtime"
	"strings"
	"text/template"

	"golang.org/x/tools/go/types"
)

const INTERFACES_TEMPLATE = `
{{range . }}{{.}}{{end}}
`
const INTERFACE_TEMPLATE = `
// {{.Name}} was extracted from {{.OriginalName}}
type {{.Name}} interface {
{{range $name,$body := .Methods}}
{{$name}}{{$body}}{{end}}
}
`

// Interfaces is a collection of Interface
type Interfaces map[string]*Interface

func (interfaces Interfaces) String() string {
	tpl, err := template.New("interfaces").Parse(INTERFACES_TEMPLATE)
	ExitOnError(err)
	b := &bytes.Buffer{}
	err = tpl.Execute(b, interfaces)
	ExitOnError(err)
	out, err := format.Source(b.Bytes())
	ExitOnError(err)
	return string(out)
}

// Interface represent an interface declaration
type Interface struct {
	Name         string
	OriginalName string
	Methods      map[string]string
}

func (i Interface) String() string {
	tpl, err := template.New("interface").Parse(INTERFACE_TEMPLATE)
	ExitOnError(err)
	b := &bytes.Buffer{}
	ExitOnError(tpl.Execute(b, i))
	return b.String()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var (
		structTypes      string
		directory        string
		structTypeFilter = map[string]bool{}
		err              error
		interfaces       = Interfaces{}
	)
	flag.ErrHelp = errors.New("extract-interface extracts interfaces from structs with methods\nUsage: extract-interface [Options]\nOptions:\n")
	flag.StringVar(&structTypes, "types", "", "\n\tType Filter , will only extract interfaces from listed types, \n\texample : -types=MyStuctType1,MyStructType2")
	flag.StringVar(&directory, "dir", "", "Directory where to find structs to extract interfaces from, \n\texample: -dir=/My/Directory")
	flag.Parse()
	if structTypes != "" {
		for _, structType := range strings.Split(strings.Trim(structTypes, " "), ",") {
			structTypeFilter[structType] = true
		}
	}
	if directory == "" {
		directory, err = os.Getwd()
		ExitOnError(err)
	} else {
		stats, err := os.Stat(directory)
		ExitOnError(err)
		if !stats.IsDir() {
			log.Fatalf("%q is not a directory", directory)
		}
	}
	fileset := token.NewFileSet()
	packages, err := parser.ParseDir(fileset, directory, nil, parser.AllErrors)
	for _, package_ := range packages {
		intrfcs := ExtractInterface(package_, structTypeFilter)
		for typeName, interface_ := range intrfcs {
			interfaces[typeName] = interface_
		}

	}
	fmt.Fprint(os.Stdout, interfaces)
}

// ExtractInterface extract interfaces from structs in a file.
// Struct types from which the interface must be extracted are filtered
// structTypeFilter map. If the length of the map is 0 , then struct types are not filtered
func ExtractInterface(Package ast.Node, structTypeFilter map[string]bool) Interfaces {
	all := false
	interfaces := Interfaces{}
	if len(structTypeFilter) == 0 {
		all = true
	}
	ast.Inspect(Package, func(node ast.Node) bool {
		if method, ok := node.(*ast.FuncDecl); ok == true && method.Recv != nil {
			typeName := strings.TrimLeft(types.ExprString(method.Recv.List[0].Type), "*")
			if !structTypeFilter[typeName] && all == false {
				return false
			}
			if interfaces[typeName] == nil {
				interfaces[typeName] = &Interface{
					OriginalName: typeName,
					Name:         typeName + "Interface",
					Methods:      map[string]string{},
				}
			}
			interfaces[typeName].Methods[method.Name.String()] = strings.TrimLeft(types.ExprString(method.Type), "func")
			return false
		}
		return true
	})
	return interfaces
}

// UTILS

func ExitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
