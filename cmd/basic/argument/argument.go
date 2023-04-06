package main

import (
	"flag"
	"fmt"
)

func main() {
	p := flag.Int("p", 1000, "port")

	flag.Parse()

	fmt.Println("Port : ", *p)
}