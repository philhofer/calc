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
	f, err := parser.ParseFile(fset, filepath.Join(thisDir, file), nil, parser.ParseComments)
	if err != nil {
		return f, err
	}
	return f, nil
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

type idMatcher struct {
	old string
	to  string
}

func (i *idMatcher) Visit(n ast.Node) ast.Visitor {
	if id, ok := n.(*ast.Ident); ok && id.Name == i.old {
		id.Name = i.to
	}
	return i
}

func SetParseDir(dir string) {
	thisDir = dir
}

func replaceIds(f *ast.File, fname string) {
	rep := idReplacer(fname)
	for _, d := range f.Decls {
		ast.Walk(rep, d)
	}
}

func WriteDerivImpl(w io.Writer, fname string, pkg string) error {
	if derivFile == nil {
		if err := parseDeriv(); err != nil {
			return err
		}
	}
	replaceIds(derivFile, fname)
	derivFile.Name.Name = pkg
	return printer.Fprint(w, fset, derivFile)
}

func WriteRootFindImpl(w io.Writer, fname string, pkg string) error {
	if rfindFile == nil {
		if err := parserFind(); err != nil {
			return err
		}
	}
	replaceIds(rfindFile, fname)
	rfindFile.Name.Name = pkg
	return printer.Fprint(w, fset, rfindFile)
}

func WriteRombergImpl(w io.Writer, fname string, pkg string) error {
	if integFile == nil {
		if err := parseInteg(); err != nil {
			return err
		}
	}
	replaceIds(integFile, fname)
	integFile.Name.Name = pkg
	if err := printer.Fprint(w, fset, integFile); err != nil {
		return err
	}
	// now we need to do implement the new 'trans' and 'transIntegral'
	// functions

	// replace find decl for __trap; replace calls to __ff
	// with calls to __ftrans, then print the node
	trapname := fname + "trap"
	for _, d := range integFile.Decls {
		if fd, ok := d.(*ast.FuncDecl); ok && fd.Name.Name == trapname {
			fd.Name.Name = fname + "transtrap"

			// replace calls to 'fname' with 'fname'+'ftrans'
			ast.Walk(&idMatcher{old: fname, to: fname + "ftrans"}, fd.Body)
			io.WriteString(w, "\n\n")
			if err := printer.Fprint(w, fset, fd); err != nil {
				return err
			}
			break
		}
	}

	// implement __transIntegral
	integname := fname + "Integral"
	for _, d := range integFile.Decls {
		if fd, ok := d.(*ast.FuncDecl); ok && fd.Name.Name == integname {
			fd.Name.Name = fname + "transIntegral"

			// replace calls to __trap with calls to __transtrap
			ast.Walk(&idMatcher{old: fname + "trap", to: fname + "transtrap"}, fd.Body)
			io.WriteString(w, "\n\n")
			if err := printer.Fprint(w, fset, fd); err != nil {
				return err
			}
			break
		}
	}
	return nil
}
