package dprint

import (
	"fmt"
	"testing"
)

type testCase struct {
	in  interface{}
	out string
}

type example0 struct{}

type example1 struct {
	first float32
}

type example2 struct {
	first  string
	second int
}

type example3 struct {
	first  int
	second int
	third  int
}

func TestDump(t *testing.T) {
	Color = false

	testCases := func(cases []testCase) {
		for _, c := range cases {
			dump := SDump(c.in)
			if fmt.Sprintf("%q", c.out) != fmt.Sprintf("%q", dump) {
				t.Errorf("\nExpected:\n%q\nGot:\n%q", c.out, dump)
			}
		}
	}

	testCases([]testCase{
		{string(`1`), `"1"`},
		{int(1), `1`},
		{int8(1), `1`},
		{int16(1), `1`},
		{int32(1), `1`},
		{int64(1), `1`},
		{float32(1.0), `1.000000`},
		{float64(1.0), `1.000000`},
		{example2{"ben", 12}, `example2{
    first: "ben",
    second: 12,
}`},
		{[]int{1, 2, 3, 4}, `[]int{
    1,
    2,
    3,
    4,
}`},
		{[]interface{}{example0{}, nil}, `[]interface {}{
    example0{
    },
    nil,
}`},
	})
}

func TestTree(t *testing.T) {
	Color = false

	testCases := func(cases []testCase) {
		for _, c := range cases {
			tree := STree(c.in)
			if fmt.Sprintf("%q", c.out) != fmt.Sprintf("%q", tree) {
				t.Errorf("\nExpected:\n%q\nGot:\n%q", c.out, tree)
			}
		}
	}

	testCases([]testCase{
		{string(`1`), `string "1"
`},
		{int(1), `int 1
`},
		{int8(1), `int8 1
`},
		{int16(1), `int16 1
`},
		{int32(1), `int32 1
`},
		{int64(1), `int64 1
`},
		{float32(1.0), `float32 1.000000
`},
		{float64(1.0), `float64 1.000000
`},
		{[]int{1, 2, 3, 4}, `[]int ─┬─ 0. int 1
       ├─ 1. int 2
       ├─ 2. int 3
       └─ 3. int 4
`},
		{[]int{1}, `[]int ─── 0. int 1
`},
		{[]int{}, `[]int
`},
		{bool(false), `bool false
`},
		{example2{"example", 123}, `dprint.example2 ─┬─ first string "example"
                 └─ second int 123
`},
		{example0{}, `dprint.example0
`},
		{example1{2.5}, `dprint.example1 ─── first float32 2.500000
`},
		{example3{1, 2, 3}, `dprint.example3 ─┬─ first int 1
                 ├─ second int 2
                 └─ third int 3
`},
		{[]interface{}{example0{}, nil}, `[]interface {} ─┬─ 0. dprint.example0
                └─ 1. nil
`},
	})

	Packages = false

	testCases([]testCase{testCase{
		example2{"example", 123}, `example2 ─┬─ first string "example"
          └─ second int 123
`},
	})

	Color = true

	testCases([]testCase{testCase{
		example2{"example2", 321}, "\x1b[34mexample2\x1b[0m ─┬─ \x1b[36mfirst\x1b[0m\x1b[34m string\x1b[0m\x1b[33m \"example2\"\x1b[0m\n          └─ \x1b[36msecond\x1b[0m\x1b[34m int\x1b[0m\x1b[32m 321\x1b[0m\n"},
	})

}
