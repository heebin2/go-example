package main

import (
	"flag"
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	efile := flag.String("ini", ".env", ".env file path")

	flag.Parse()

	if err := godotenv.Load(*efile); err != nil {
		fmt.Println(err)
	}
}