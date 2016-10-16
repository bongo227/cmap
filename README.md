#dprint - Deep pretty printer in color#
[![codecov](https://codecov.io/gh/bongo227/dprint/branch/master/graph/badge.svg)](https://codecov.io/gh/bongo227/dprint)
[![Build Status](https://travis-ci.org/bongo227/dprint.svg?branch=master)](https://travis-ci.org/bongo227/dprint)
[![](https://godoc.org/github.com/bongo227/dprint?status.svg)](http://godoc.org/github.com/bongo227/dprint)
##Usage##
####Tree####
```go
// Pretty print a data structure (as a tree)
dprint.Tree(interface{})

// Pretty printed string of a data structure (as a tree)
dprint.STree(interface{})
```
####Dump####
```go
// Pretty print a data structure (in a go format)
dprint.Dump(interface{})

// Pretty printed string of a data structure (in a go format)
dprint.SDump(interface{})
```