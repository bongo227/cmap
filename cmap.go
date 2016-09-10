package cmap

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func Dump(item interface{}, name string) {
	itype := reflect.TypeOf(item)
	ivalue := reflect.ValueOf(item)

	printThis(name, itype, ivalue, "")
}

func printThis(name string, itype reflect.Type, ivalue reflect.Value, depth string) {
	cname := color.New(color.FgHiCyan).SprintfFunc()
	ctype := color.New(color.FgBlue).SprintfFunc()
	cstring := color.New(color.FgYellow).SprintfFunc()
	cint := color.New(color.FgHiGreen).SprintfFunc()
	cposition := color.New(color.FgHiBlack).SprintfFunc()
	cnil := color.New(color.FgRed).SprintfFunc()

	typeName := strings.Replace(itype.String(), "compiler.", "", 100)

	if itype.Kind() != reflect.Interface {
		s := ""
		if name != "" {
			s = cname("%s", name) + ctype(" %s", typeName)
			depth += strings.Repeat(" ", len(name)+1+len(typeName))
		} else {
			s = ctype("%s", typeName)
			depth += strings.Repeat(" ", len(typeName))
		}

		fmt.Print(s)
	}

	if itype.Kind() == reflect.Slice {

		if ivalue.Len() == 0 {
			fmt.Println()
		}

		for i := 0; i < ivalue.Len(); i++ {
			if ivalue.Len() == 1 {
				fmt.Printf(" ─── %s ", cposition("%d.", i))
				printThis("", ivalue.Index(i).Type(), ivalue.Index(i), depth+"        ")
			} else if i == 0 {
				fmt.Printf(" ─┬─ %s ", cposition("%d.", i))
				printThis("", ivalue.Index(i).Type(), ivalue.Index(i), depth+"  │     ")
			} else if i == ivalue.Len()-1 {
				fmt.Printf("%s  └─ %s ", depth, cposition("%d.", i))
				printThis("", ivalue.Index(i).Type(), ivalue.Index(i), depth+"        ")
			} else {
				fmt.Printf("%s  ├─ %s ", depth, cposition("%d.", i))
				printThis("", ivalue.Index(i).Type(), ivalue.Index(i), depth+"  │     ")
			}
		}
	} else if itype.Kind() == reflect.String {
		s := cstring(" %s", strconv.Quote(ivalue.String()))
		fmt.Printf("%s\n", s)
	} else if itype.Kind() == reflect.Int {
		s := cint(" %d", ivalue.Int())
		fmt.Printf("%s\n", s)
	} else if itype.Kind() == reflect.Struct {
		if ivalue.NumField() == 0 {
			fmt.Println()
		}

		for i := 0; i < ivalue.NumField(); i++ {
			newDepth := ""
			if ivalue.NumField() == 1 {
				fmt.Printf(" ─── ")
				newDepth = depth + "     "
			} else if i == 0 {
				fmt.Printf(" ─┬─ ")
				newDepth = depth + "  |  "
			} else if i == ivalue.NumField()-1 {
				fmt.Printf("%s  └─ ", depth)
				newDepth = depth + "     "
			} else {
				fmt.Printf("%s  ├─ ", depth)
				newDepth = depth + "  |  "
			}

			printThis(itype.Field(i).Name, ivalue.Field(i).Type(), ivalue.Field(i), newDepth)
		}
	} else if itype.Kind() == reflect.Interface {
		if ivalue.Elem() == reflect.ValueOf(nil) {
			fmt.Printf(cnil("nil\n"))
		} else {
			printThis(name, ivalue.Elem().Type(), ivalue.Elem(), depth)
		}
	} else {
		panic("Unrecognized type")
	}
}
