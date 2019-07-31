package main

import (
	"custom/chapter1"
	"flag"
	"fmt"
	"strings"
)
type Fahrenheit float64
var n = flag.Bool("n", false, "omit trailing newline")
var sep = flag.String("s", "/", "separator")
func main() {
	echo4()
	p := new(int)
	fmt.Println(p)
	medals := []string{"gold", "silver"}
	fmt.Println(len(medals))
	fmt.Println(medals)

	fa := Fahrenheit(float64(2.02))
	fmt.Println(fa)
	fmt.Println(fa)
	fmt.Println(fa + 2.02)

	// not support fmt.Println(fa == float64(2.02))
	chapter1.ControlFlow()

}
func echo4() {
	flag.Parse()
	fmt.Print(strings.Join(flag.Args(), *sep))
	if !*n {
		fmt.Println()
	}

}