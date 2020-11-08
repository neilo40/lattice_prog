package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lattice_prog/internal/app/prog"
)

var binFile string

func initOpts() {
	flag.StringVar(&binFile, "binfile", "", "Path to binary file")
}

func checkOpts() {
	binFileFlag := flag.Lookup("binfile")
	if binFileFlag.Value.String() == "" {
		fmt.Println("binfile is a mandatory option")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	initOpts()
	flag.Parse()
	checkOpts()
	prog.Init()
	prog.Program(binFile)
}
