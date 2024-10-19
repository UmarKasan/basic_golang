package main

import (
	"fmt"
)

func main() {
	var a int    // if not set, will auto assign as 0, auto assign as int32
	var b int8   // 00000010 = 2 -127 to 127
	var c uint8  // 00000010 = 2 0 to 255
	var d uint32 // 00000010 = 2 0 to 4294967295

	var e float32 // 4000.0004
	var f bool    // true or false

	// Multiline comment
	/* var g, h, i, j, k, l, m, n, o, p, q, r,
	s, t, u, v, w, x, y, z int */

	a = -2000000
	b = 127
	c = 255
	d = 4294967295

	e = 4000.0004

	// Addition
	fmt.Println(b + int8(e))
	fmt.Println(float32(b) + e) // notice the value is wrong. happens with floats

	// Console print
	fmt.Println("Hello World") // prints hello world

	fmt.Println(a, b, c, d, e, f) // have to utilise, unable to run due to error

	// Print each variable individually
	fmt.Println("\nPrinting each variable:")
	variables := []interface{}{a, b, c, d, e, f}
	varNames := []string{"a", "b", "c", "d", "e", "f"}
	for i, v := range variables {
		fmt.Printf("%s: %v\n", varNames[i], v)
	}

	// Using for loop for variable from (A-Z):
	for i := 65; i <= 90; i++ { // A = 65 z = 90
		fmt.Printf("%c: %d\n", rune(i), i) // rune() converts int to ASCII char, using int() returns ASCII to int
	}

	// Using for loop for variable from (A-Z):
	for i, v := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		fmt.Printf("%c = %d\n", rune(v), i)
	}

}
