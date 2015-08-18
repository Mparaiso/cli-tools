// getset generates getters and setters for your structs

package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/interactiv/array"
	"golang.org/x/tools/go/types"
	"golang.org/x/tools/imports"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	wd, err := os.Getwd()
	exitOnError(err, "error getting working directory")
	log := log.New(os.Stdout, "", 0)
	Type := flag.String("type", "", "type from which getters and setters will be generated, declare multiple types by separating them with a comma, for instance: -type=Foo,Bar,Baz")
	Dir := flag.String("dir", wd, "package directory.defaults to working directory")
	File := flag.Bool("file", false, "write to file, true of false, false by default, if false then will write to the standard out stdout.")
	flag.Parse()
	if *Type == "" {
		flag.Usage()
		os.Exit(1)
	}
	*Type = strings.Replace(*Type, " ", "", -1)
	typeList := strings.Split(*Type, ",")
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, *Dir, nil, 0)
	exitOnError(err, "error parsing directory")

	for _, pkg := range pkgs {
		header := "\n// this file was generated by getset tool \n// github.com/interactiv/getset\n// " + time.Now().String() + " \n\n\n package " + pkg.Name
		sourceStringWithPackage := []string{header}
		ast.Inspect(pkg, func(n ast.Node) bool {
			if typeSpec, ok := n.(*ast.TypeSpec); ok != true {
				return true
			} else if array.IndexOf(typeList, typeSpec.Name.String()) < 0 && *Type != "*" {
				return true
			} else if structType, ok := typeSpec.Type.(*ast.StructType); ok != true {
				return true
			} else if len(structType.Fields.List) == 0 {
				return true
			} else {
				sourceStringWithPackage = append(sourceStringWithPackage, GenerateGettersAndSettersForStructType(typeSpec.Name.String(), structType))
			}
			return false
		})
		if len(sourceStringWithPackage) == 1 {
			continue // no structs or fields in struct found, inspect the next package
		}
		// string from []string
		reducedSource := array.Reduce(sourceStringWithPackage, func(result, elem interface{}, i int) interface{} {
			return result.(string) + elem.(string)
		}, "").(string)
		// fix imports and format
		source, err := imports.Process("getters_and_setters.go", []byte(reducedSource), nil)
		exitOnError(err, "error fixing imports and formating source")
		// render file if *File or log to std
		if !*File {

			log.Println(string(source))
		} else {
			file, err := os.OpenFile(path.Join(*Dir, pkg.Name+"_getters_and_setters.go"), os.O_CREATE, os.ModePerm)
			exitOnError(err, "error opening file")
			_, err = file.WriteString(string(source))
			exitOnError(err, "error writing file")
		}

	}

}

func GenerateGettersAndSettersForStructType(typeName string, structType *ast.StructType) (result string) {
	result += "\n/*\n * Getters and setters for struct type " + typeName + "\n*/\n\n"
	for _, field := range structType.Fields.List {
		fieldTypeString := types.ExprString(field.Type)
		result += GenerateGetterAndSetterString(typeName, field.Names[0].String(), fieldTypeString)
	}
	return
}

// exitOnError logs an error and exits the program if error is not nil
func exitOnError(Error error, extra ...string) {
	if Error != nil {
		log.Fatal(fmt.Sprintln(Error, extra))
		debug.PrintStack()
	}
}

// GenerateGetterAndSetterString returns both the result of GenerateGetterString and GenerateSetterString as a string
func GenerateGetterAndSetterString(typeName, propertyName, propertyType string) string {
	return GenerateGetterString(typeName, propertyName, propertyType) + " " + GenerateSetterString(typeName, propertyName, propertyType)
}

// GenerateGetterString returns a string representing a getter function for field propertyName with a receiver of type typeName
func GenerateGetterString(typeName, propertyName, propertyType string) string {
	propWithCaps := CapitalCase(propertyName)
	smallCasedTypeName := SmallCase(typeName)
	return sprintf(get, typeName, smallCasedTypeName, propertyName, propWithCaps, propertyType)
}

// GenerateSetterString returns a string representing a setter function for field propertyName with a receiver of type typeName
func GenerateSetterString(typeName, propertyName, propertyType string) string {
	propWithCaps := CapitalCase(propertyName)
	smallCasedTypeName := SmallCase(typeName)
	return sprintf(set, typeName, smallCasedTypeName, propertyName, propWithCaps, propertyType)
}

var (
	print   = fmt.Println
	printf  = fmt.Printf
	sprintf = fmt.Sprintf
	get     = `
// Get%[4]s returns a %[5]s
func (%[2]s %[1]s) %[4]s()%[5]s{
    return %[2]s.%[3]s
}
`
	set = `
// Set%[4]s sets *%[1]s.%[2]s and returns *%[1]s
func (%[2]s *%[1]s) Set%[4]s(%[3]s %[5]s)*%[1]s{
    %[2]s.%[3]s = %[3]s
    return %[2]s
}
`
)

// CapitalCase returns a capital cased string
func CapitalCase(word string) (capitalCasedWord string) {
	switch len(word) {
	case 0:
		capitalCasedWord = ""
	case 1:
		capitalCasedWord = strings.ToUpper(word)
	default:
		firstLetter := strings.ToUpper(string(word[0]))
		capitalCasedWord = firstLetter + word[1:]
	}
	return
}

// SmallCase returns a small cased string
func SmallCase(word string) (smallCasedWorld string) {
	switch len(word) {
	case 0:
		smallCasedWorld = ""
	case 1:
		smallCasedWorld = strings.ToLower(word)
	default:
		firstLetter := strings.ToLower(string(word[0]))
		smallCasedWorld = firstLetter + word[1:]
	}
	return
}
