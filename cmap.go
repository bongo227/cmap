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

	out := ""

	printThis(name, itype, ivalue, "", &out)
	fmt.Println(out)
}

func SDump(item interface{}, name string) string {
	itype := reflect.TypeOf(item)
	ivalue := reflect.ValueOf(item)

	out := ""

	printThis(name, itype, ivalue, "", &out)
	return out
}

func printThis(name string, itype reflect.Type, ivalue reflect.Value, depth string, out *string) {

	cname := color.New(color.FgHiCyan).SprintfFunc()
	ctype := color.New(color.FgBlue).SprintfFunc()
	cstring := color.New(color.FgYellow).SprintfFunc()
	cint := color.New(color.FgHiGreen).SprintfFunc()
	cbool := color.New(color.FgHiMagenta).SprintfFunc()
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

		*out += s
	}

	if itype.Kind() == reflect.Slice {

		if ivalue.Len() == 0 {
			*out += "\n"
		}

		for i := 0; i < ivalue.Len(); i++ {
			if ivalue.Len() == 1 {
				*out += fmt.Sprintf(" ─── %s ", cposition("%d.", i))
				printThis("", ivalue.Index(i).Type(), ivalue.Index(i), depth+"        ", out)
			} else if i == 0 {
				*out += fmt.Sprintf(" ─┬─ %s ", cposition("%d.", i))
				printThis("", ivalue.Index(i).Type(), ivalue.Index(i), depth+"  │     ", out)
			} else if i == ivalue.Len()-1 {
				*out += fmt.Sprintf("%s  └─ %s ", depth, cposition("%d.", i))
				printThis("", ivalue.Index(i).Type(), ivalue.Index(i), depth+"        ", out)
			} else {
				*out += fmt.Sprintf("%s  ├─ %s ", depth, cposition("%d.", i))
				printThis("", ivalue.Index(i).Type(), ivalue.Index(i), depth+"  │     ", out)
			}
		}
	} else if itype.Kind() == reflect.String {
		s := cstring(" %s", strconv.Quote(ivalue.String()))
		*out += fmt.Sprintf("%s\n", s)
	} else if itype.Kind() == reflect.Int {
		s := cint(" %d", ivalue.Int())
		*out += fmt.Sprintf("%s\n", s)
	} else if itype.Kind() == reflect.Float32 {
		s := cint(" %f", ivalue.Float())
		*out += fmt.Sprintf("%s\n", s)
	} else if itype.Kind() == reflect.Float64 {
		s := cint(" %f", ivalue.Float())
		*out += fmt.Sprintf("%s\n", s)
	} else if itype.Kind() == reflect.Bool {
		s := cbool(" %t", ivalue.Bool())
		*out += fmt.Sprintf("%s\n", s)
	} else if itype.Kind() == reflect.Struct {
		if ivalue.NumField() == 0 {
			*out += "\n"
		}

		for i := 0; i < ivalue.NumField(); i++ {
			newDepth := ""
			if ivalue.NumField() == 1 {
				*out += " ─── "
				newDepth = depth + "     "
			} else if i == 0 {
				*out += " ─┬─ "
				newDepth = depth + "  |  "
			} else if i == ivalue.NumField()-1 {
				*out += fmt.Sprintf("%s  └─ ", depth)
				newDepth = depth + "     "
			} else {
				*out += fmt.Sprintf("%s  ├─ ", depth)
				newDepth = depth + "  |  "
			}

			printThis(itype.Field(i).Name, ivalue.Field(i).Type(), ivalue.Field(i), newDepth, out)
		}
	} else if itype.Kind() == reflect.Interface {
		if ivalue.Elem() == reflect.ValueOf(nil) {
			*out += cnil("nil\n")
		} else {
			printThis(name, ivalue.Elem().Type(), ivalue.Elem(), depth, out)
		}
	} else {
		*out += fmt.Sprintf("%s\n", "Unrecognized type")
	}
}
