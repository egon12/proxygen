package main

import (
	"os"

	"github.com/egon12/proxygen"
)

func main() {
	filename := os.Args[1]
	names := os.Args[2]

	err := proxygen.Generate(filename, names)
	if err != nil {
		panic(err)
	}
}
