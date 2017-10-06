package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

type parsedData struct {
	Filename           string
	ResourceName       string
	ResourceDefinition []docDefinition
}

func parseProvider(fset *token.FileSet, fp string) parsedData {
	// Perform the parsing of the Go code
	// or fail if we couldn't parse it (this can be anything from
	// a mistype to the file not being there)
	parsed, err := parser.ParseFile(fset, fp, nil, 0)
	if err != nil {
		errexit("Unable to parse file %q: %s", filepath.Base(fp), err.Error())
	}

	// Call the walker by creating a reference first
	// so we can get the list of attributes after
	walker := new(visitorProvider)
	ast.Walk(walker, parsed)

	// Generate a structure with the proper details
	return parsedData{
		Filename:           filepath.Base(fp),
		ResourceName:       "nn",
		ResourceDefinition: walker.results,
	}
}
