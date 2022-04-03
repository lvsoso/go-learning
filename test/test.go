package main

// #include "test.h"
// #cgo LDFLAGS: libtest.a -lstdc++
import "C"

func main() {
	C.test()
}
