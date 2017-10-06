package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

type parsedData struct {
	Filename           string
	ResourceName       string
	IsMainProvider     bool
	ResourceDefinition []docDefinition
}

var rscCleaner = strings.NewReplacer(
	"resource_", "",
	".go", "",
)

func parseProvider(providerName string, fset *token.FileSet, fp string) parsedData {
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

	// Generate a name for the file
	filename := filepath.Base(fp)
	resourcename := ""

	// If it's a resource, generate the proper resource name
	if strings.HasPrefix(filename, "resource") {
		// This will generate something like "provider_resource_name"
		resourcename = providerName + "_" + rscCleaner.Replace(filename)
	}

	// If the filename is the provider itself, we change the resource name
	// to just be the provider name
	if filename == "provider.go" {
		resourcename = providerName
	}

	// Generate a structure with the proper details
	return parsedData{
		Filename:           filename,
		ResourceName:       resourcename,
		IsMainProvider:     resourcename == providerName,
		ResourceDefinition: walker.results,
	}
}
