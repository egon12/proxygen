package main

import (
	"log"
	"os"

	"github.com/egon12/proxygen"
)

func main() {
	filename := os.Args[1]
	names := os.Args[2]

	err := proxygen.Generate(filename, names)
	if err != nil {
		log.Println(err)
	}
}
