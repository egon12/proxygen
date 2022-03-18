package main

import (
	"flag"

	"github.com/egon12/proxygen"
)

func main() {

	var filename, names, outputDir, packageName string

	flag.StringVar(&filename, "filename", "", "file to be process")
	flag.StringVar(&names, "names", "", "interface name that need to be process")
	flag.StringVar(&outputDir, "output", ".", "folder output")
	flag.StringVar(&packageName, "package", "", "package name (default use interface package")

	flag.Parse()

	if filename == "" || names == "" {
		flag.Usage()
		return
	}

	err := proxygen.Generate(filename, names, outputDir, packageName)
	if err != nil {
		panic(err)
	}
}
