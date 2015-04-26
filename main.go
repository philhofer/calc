package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/philhofer/calc/pgen"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	thisDir string

	fname  = flag.String("func", "", "function name")
	method = flag.String("method", "", "method to generate")
	gopkg  = os.Getenv("GOPACKAGE")
	gofile = os.Getenv("GOFILE")
)

func init() {
	_, thisfile, _, _ := runtime.Caller(0)
	thisDir = filepath.Dir(thisfile)
}

func fail(s string) {
	fmt.Println("Error:", s)
	os.Exit(1)
}

func main() {
	flag.Parse()
	pgen.SetParseDir(thisDir)
	newfile := strings.TrimSuffix(gofile, ".go") + "_" + *method + ".go"
	file, err := os.Create(newfile)
	if err != nil {
		fail(err.Error())
	}
	defer file.Close()
	b := bufio.NewWriter(file)
	switch *method {
	case "deriv":
		err = pgen.WriteDerivImpl(b, *fname, gopkg)
	case "root":
		err = pgen.WriteRootFindImpl(b, *fname, gopkg)
	default:
		err = fmt.Errorf("no method named %q available to generate", *method)
	}
	if err == nil {
		err = b.Flush()
	}
	if err != nil {
		file.Close()
		os.Remove(newfile)
		fail(err.Error())
	}
}
