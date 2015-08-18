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
	"go/parser"
	"go/token"
	"log"
	"testing"

	"github.com/interactiv/expect"
)

func TestExtractInterface(t *testing.T) {
	e := expect.New(t)
	e.Expect(func() {
		fileset := token.NewFileSet()
		pkg, err := parser.ParseFile(fileset, "foo.go", TEST_PACKAGE, parser.AllErrors)
		if err != nil {
			log.Fatal(err)
		}
		interfaces := ExtractInterface(pkg, map[string]bool{})
		e.Expect(interfaces.String()).ToContain(EXPECTED_RESULT)
	}).Not().ToPanic()
}

const TEST_PACKAGE = `
package foo

type NewString string

type Foo struct{}

func (f Foo)Do()String{
    return  "doing"
}

type Bar struct {
	i      int
	things []string
}
func DoSomething(s string){}
func (b Bar) ReturnString() string {}
func (b *Bar) SetStuff(i int) *Bar {
	b.i = i
	return b
}
func (b Bar) AddThings(things ...string) *Bar {
	b.things = things
	return b
}
`

const EXPECTED_RESULT = `

// BarInterface was extracted from Bar
type BarInterface interface {
	AddThings(things ...string) *Bar
	ReturnString() string
	SetStuff(i int) *Bar
}

// FooInterface was extracted from Foo
type FooInterface interface {
	Do() String
}

`
