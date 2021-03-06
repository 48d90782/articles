package main

import (
	"C"
)

// we need to have empty main in package main :)
// because -buildmode=c-shared requires exactly one main package
func main() {

}

//export Fib
func Fib(n C.int) C.int {
	if n == 1 || n == 2 {
		return 1
	}

	return Fib(n-1) + Fib(n-2)
}
