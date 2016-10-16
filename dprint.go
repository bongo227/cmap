package dprint

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

const (
	escape = "\x1b"
	reset  = escape + "[0m"
)

var (
	cposition = colorFunc("30") // Black/Grey
	cnil      = colorFunc("31") // Red
	cint      = colorFunc("32") // Green
	cstring   = colorFunc("33") // Yellow
	ctype     = colorFunc("34") // Blue
	cbool     = colorFunc("35") // Magenta
	cname     = colorFunc("36") // Cyan

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

// Dump prints data structures in color in a golang style
func Dump(item interface{}) {
	fmt.Println(SDump(item))
}

// SDump returns the string result of Dump
func SDump(item interface{}) string {
	s := ""
	itype := reflect.TypeOf(item)
	ivalue := reflect.ValueOf(item)

	linePrint(itype, ivalue, "", &s)
	return s
}

func linePrint(itype reflect.Type, ivalue reflect.Value, depth string, out *string) {
	typeName := itype.String()

	// Remove local package
	if lastDot := strings.LastIndex(typeName, "."); lastDot > 0 && unicode.IsLower(rune(typeName[lastDot+1])) {
		typeName = typeName[lastDot+1 : len(typeName)]
	}

	switch itype.Kind() {

	// Strings
	case reflect.String:
		*out += cstring("%q", ivalue.String())

	// Intergers
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		*out += cint("%d", ivalue.Int())

	// Floats
	case reflect.Float32, reflect.Float64:
		*out += cint("%f", ivalue.Float())

	// Bools
	case reflect.Bool:
		*out += cbool("%t", ivalue.Bool())

	// Structures
	case reflect.Struct:
		*out += fmt.Sprintf("%s{\n", ctype(typeName))

		for i := 0; i < ivalue.NumField(); i++ {
			*out += fmt.Sprintf("%s%s: ", depth+"    ", cname(itype.Field(i).Name))
			linePrint(ivalue.Field(i).Type(), ivalue.Field(i), depth+"    ", out)
			*out += fmt.Sprintf(",\n")
		}

		*out += fmt.Sprintf("%s}", depth)

	// Slices
	case reflect.Slice:
		*out += fmt.Sprintf("[]%s{\n", ctype(typeName[2:]))

		for i := 0; i < ivalue.Len(); i++ {
			*out += fmt.Sprintf("%s", depth+"    ")
			linePrint(ivalue.Index(i).Type(), ivalue.Index(i), depth+"    ", out)
			*out += fmt.Sprintf(",\n")
		}

		*out += fmt.Sprintf("%s}", depth)

	default:
		*out += "Unrecognised type"
	}
}

// Tree pretty prints data structures in color in the style of horizontal tree
func Tree(item interface{}) {
	fmt.Println(STree(item))
}

// STree returns the string result of tree
func STree(item interface{}) string {
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
