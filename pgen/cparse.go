package pgen

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"path/filepath"
	"strings"
)

const (
	derivName = "calc/deriv.go"
	integName = "calc/romberg.go"
	rfindName = "calc/root.go"
)

var (
	derivFile *ast.File
	integFile *ast.File
	rfindFile *ast.File
	fset      = token.NewFileSet()
	thisDir   string
)

func parseFile(thisDir string, file string) (*ast.File, error) {
	return parser.ParseFile(fset, filepath.Join(thisDir, file), nil, parser.ParseComments)
}

func parseDeriv() error {
	var err error
	derivFile, err = parseFile(thisDir, derivName)
	return err
}

func parseInteg() error {
	var err error
	integFile, err = parseFile(thisDir, integName)
	return err
}

func parserFind() error {
	var err error
	rfindFile, err = parseFile(thisDir, rfindName)
	return err
}

type idReplacer string

// replace all ids with '__' prefixes to their proper names
func (i idReplacer) Visit(n ast.Node) ast.Visitor {
	if id, ok := n.(*ast.Ident); ok {
		if strings.HasPrefix(id.Name, "__") {
			if id.Name == "__ff" {
				id.Name = string(i)
			} else {
				id.Name = string(i) + id.Name[2:]
			}
		}
	}
	return i
}

func SetParseDir(dir string) {
	thisDir = dir
}

func WriteDerivImpl(w io.Writer, fname string, pkg string) error {
	if derivFile == nil {
		if err := parseDeriv(); err != nil {
			return err
		}
	}
	rep := idReplacer(fname)
	for _, d := range derivFile.Decls {
		ast.Walk(rep, d)
	}
	derivFile.Name.Name = pkg
	return printer.Fprint(w, fset, derivFile)
}

func WriteRootFindImpl(w io.Writer, fname string, pkg string) error {
	if rfindFile == nil {
		if err := parserFind(); err != nil {
			return err
		}
	}
	rep := idReplacer(fname)
	for _, d := range rfindFile.Decls {
		ast.Walk(rep, d)
	}
	rfindFile.Name.Name = pkg
	return printer.Fprint(w, fset, rfindFile)
}
