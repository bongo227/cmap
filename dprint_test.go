package dprint

import (
	"fmt"
	"testing"
)

type testCase struct {
	in  interface{}
	out string
}

type example1 struct {
	first  string
	second int
}

type example2 struct{}

type example3 struct {
	first float32
}

type example4 struct {
	first  int
	second int
	third  int
	forth  int
}

func TestDump(t *testing.T) {
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
		{example1{"example", 123}, `dprint.example1 ─┬─ first string "example"
                 └─ second int 123
`},
		{example2{}, `dprint.example2
`},
		{example3{2.5}, `dprint.example3 ─── first float32 2.500000
`},
		{example4{1, 2, 3, 4}, `dprint.example4 ─┬─ first int 1
                 ├─ second int 2
                 ├─ third int 3
                 └─ forth int 4
`},
		{[]interface{}{example2{}, nil}, `[]interface {} ─┬─ 0. dprint.example2
                └─ 1. nil
`},
	})

	Packages = false

	testCases([]testCase{testCase{
		example1{"example", 123}, `example1 ─┬─ first string "example"
          └─ second int 123
`},
	})

	Color = true

	testCases([]testCase{testCase{
		example1{"example2", 321}, "\x1b[34mexample1\x1b[0m ─┬─ \x1b[36mfirst\x1b[0m\x1b[34m string\x1b[0m\x1b[33m \"example2\"\x1b[0m\n          └─ \x1b[36msecond\x1b[0m\x1b[34m int\x1b[0m\x1b[32m 321\x1b[0m\n"},
	})

}
