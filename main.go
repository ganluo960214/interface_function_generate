package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/ganluo960214/ast_extend"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
)

/*
flags
*/
var (
	flags = struct {
		Func string `validate:"required"`
		Type string `validate:"required"`
	}{}
)

func init() {
	flag.StringVar(&flags.Func, "func", "", "function name; must be set")
	flag.StringVar(&flags.Type, "type", "", "data type split by ','; must be set")

	flag.Parse()
	if err := validate.Struct(flags); err != nil {
		log.Fatalln(err)
	}
}

/*
env
*/
var (
	envs = struct {
		GoPackage string `validate:"required"`
		GoFile    string `validate:"required,check_file_exists"`
	}{
		GoPackage: os.Getenv("GOPACKAGE"),
		GoFile:    os.Getenv("GOFILE"),
	}
)

func init() {
	if err := validate.Struct(flags); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	// ast file
	fSet := token.NewFileSet()
	parseFile, err := parser.ParseFile(fSet, envs.GoFile, nil, parser.ParseComments)
	if err != nil {
		log.Fatalln(err)
	}

	// make a reader at for source file
	var rat io.ReaderAt
	bs, err := ioutil.ReadFile(envs.GoFile)
	if err != nil {
		log.Fatalln(err)
	}
	rat = bytes.NewReader(bs)

	// find func
	decl, ok := ast_extend.AstExtFindFuncByName(parseFile, "", flags.Func)
	if ok == false {
		log.Fatalln(fmt.Sprintf("not found function(%s)", flags.Func))
	}

	// find interface params
	funcs := []string{}
	for _, Type := range strings.Split(flags.Type, ",") {
		b := strings.Builder{}

		s1, err := IoReadAtByTokenPos(
			rat,
			decl.Pos(), decl.Type.Params.Pos())
		if err != nil {
			log.Fatalln(err)
		}
		b.WriteString(fmt.Sprintf("%s_%s(", s1, Type))

		for count, f := range decl.Type.Params.List {
			begin, end := f.Pos(), f.End()
			_, isInterfaceType := f.Type.(*ast.InterfaceType)
			if isInterfaceType {
				begin, end = f.Names[0].Pos(), f.Names[len(f.Names)-1].End()
			}

			p, err := IoReadAtByTokenPos(
				rat,
				begin, end)
			if err != nil {
				log.Fatalln(err)
			}
			b.WriteString(p)

			if isInterfaceType {
				b.WriteString(" " + Type)
			}

			if count+1 != len(decl.Type.Params.List) {
				b.WriteRune(',')
			}
		}
		b.WriteRune(')')

		s2, err := IoReadAtByTokenPos(
			rat,
			decl.Type.Params.End(), decl.End())
		if err != nil {
			log.Fatalln(err)
		}
		b.WriteString(s2)

		funcs = append(funcs, b.String())
	}

	t := FileTemplateContent{
		Flags:   strings.Join(os.Args, " "),
		Package: envs.GoPackage,
		Funcs:   funcs,
	}
	fc, err := t.generateContent()
	if err != nil {
		log.Fatalln(err)
	}

	formattedFC, err := format.Source(fc)
	if err != nil {
		log.Fatalln(err)
	}

	goend := ".go"
	goTestEnd := "_test.go"
	if len(envs.GoFile) < len(goend) || reflect.DeepEqual(envs.GoFile[len(envs.GoFile)-3:], goend) == false {
		log.Fatalln("file not end with '.go'")
	}
	generateFileName := envs.GoFile[:len(envs.GoFile)-len(goend)] + "_interface_function_generate" + envs.GoFile[len(envs.GoFile)-len(goend):]
	if len(envs.GoFile) > len(goTestEnd) && reflect.DeepEqual(envs.GoFile[len(envs.GoFile)-len(goTestEnd):], goTestEnd) {
		generateFileName = envs.GoFile[:len(envs.GoFile)-len(goend)] + "_interface_function_generate_test" + envs.GoFile[len(envs.GoFile)-len(goend):]
	}

	if err := ioutil.WriteFile(generateFileName, formattedFC, 0644); err != nil {
		log.Fatalln(err)
	}
}
