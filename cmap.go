package cmap

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

const (
	escape = "\x1b"
	reset  = escape + "[0m"
)

var (
	cposition = colorFunc("30")
	cnil      = colorFunc("31")
	cint      = colorFunc("32")
	cstring   = colorFunc("33")
	ctype     = colorFunc("34")
	cbool     = colorFunc("35")
	cname     = colorFunc("36")

	// Packages controls wether the package a struct comes from is printed
	Packages = true

	// Color controls wether color codes are outputed
	Color = true
)

func colorFunc(colorCode string) func(string, ...interface{}) string {
	return func(format string, a ...interface{}) string {
		if Color {
			return escape + "[" + colorCode + "m" + fmt.Sprintf(format, a...) + reset
		}

		return fmt.Sprintf(format, a...)
	}
}

// Tree pretty prints data structures in color in the style of horizontal tree
func Tree(item interface{}) {
	fmt.Println(STree(item))
}

// STree returns the string result of tree
func STree(item interface{}) string {
	if !Color {
		color.NoColor = true
		defer func() { color.NoColor = false }()
	}

	itype := reflect.TypeOf(item)
	ivalue := reflect.ValueOf(item)

	out := ""

	branchPrint("", itype, ivalue, "", &out)
	return out
}

func branchPrint(name string, itype reflect.Type, ivalue reflect.Value, depth string, out *string) {
	typeName := itype.String()
	if !Packages {
		if lastDot := strings.LastIndex(typeName, "."); lastDot > 0 {
			typeName = typeName[lastDot+1 : len(typeName)]
		}
	}

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

	switch itype.Kind() {

	// Slices
	case reflect.Slice:
		if ivalue.Len() == 0 {
			*out += "\n"
		}

		for i := 0; i < ivalue.Len(); i++ {
			if ivalue.Len() == 1 {
				*out += fmt.Sprintf(" ─── %s ", cposition("%d.", i))
				branchPrint("", ivalue.Index(i).Type(), ivalue.Index(i), depth+"        ", out)
			} else if i == 0 {
				*out += fmt.Sprintf(" ─┬─ %s ", cposition("%d.", i))
				branchPrint("", ivalue.Index(i).Type(), ivalue.Index(i), depth+"  │     ", out)
			} else if i == ivalue.Len()-1 {
				*out += fmt.Sprintf("%s  └─ %s ", depth, cposition("%d.", i))
				branchPrint("", ivalue.Index(i).Type(), ivalue.Index(i), depth+"        ", out)
			} else {
				*out += fmt.Sprintf("%s  ├─ %s ", depth, cposition("%d.", i))
				branchPrint("", ivalue.Index(i).Type(), ivalue.Index(i), depth+"  │     ", out)
			}
		}

	// Strings
	case reflect.String:
		s := cstring(" %s", strconv.Quote(ivalue.String()))
		*out += fmt.Sprintf("%s\n", s)

	// Intergers
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s := cint(" %d", ivalue.Int())
		*out += fmt.Sprintf("%s\n", s)

	// Floats
	case reflect.Float32, reflect.Float64:
		s := cint(" %f", ivalue.Float())
		*out += fmt.Sprintf("%s\n", s)

	// Booleans
	case reflect.Bool:
		s := cbool(" %t", ivalue.Bool())
		*out += fmt.Sprintf("%s\n", s)

	// Structs
	case reflect.Struct:
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

			branchPrint(itype.Field(i).Name, ivalue.Field(i).Type(), ivalue.Field(i), newDepth, out)
		}

	// Interfaces
	case reflect.Interface:
		if ivalue.Elem() == reflect.ValueOf(nil) {
			*out += cnil("nil\n")
		} else {
			branchPrint(name, ivalue.Elem().Type(), ivalue.Elem(), depth, out)
		}

	default:
		*out += fmt.Sprintf("%s\n", " Unrecognized type")

	}
}
