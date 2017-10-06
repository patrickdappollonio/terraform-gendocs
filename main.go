package main

import (
	"fmt"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Find if a second argument is passed
	if len(os.Args) != 2 {
		errexit("Usage: terraform-gendocs {go-import-path}")
	}

	// Find if there's a $GOPATH
	var gopath string
	if v := os.Getenv("GOPATH"); v != "" {
		gopath = v
	}

	// Check if we could get a $GOPATH
	if gopath == "" {
		errexit("Unable to find a $GOPATH. Make sure you set one.")
	}

	// Generate a full path based on the $GOPATH defined
	fullpath := filepath.Join(gopath, "src", os.Args[1])

	// Check if the folder requested exists
	if !exists(fullpath) {
		errexit("The path %q does not exists", fullpath)
	}

	// Get the name of the provider
	provider := filepath.Base(fullpath)
	internalName := strings.TrimPrefix(provider, "terraform-provider-")
	internalPath := filepath.Join(fullpath, internalName)

	// Find all files / folders here
	files, err := ioutil.ReadDir(internalPath)

	if err != nil {
		errexit("Unable to read folder to internal path name %q at %q. This project needs the provider to be configured as expected by Terraform", internalName, internalPath)
	}

	// Create a list of all the important files
	tffiles := make(map[string]string)

	// Iterate over all found files and add those we think are the useful ones
	for _, v := range files {
		// Find if the file is the provider configuration
		if v.Name() == "provider.go" {
			tffiles[v.Name()] = filepath.Join(internalPath, v.Name())
		}

		// Find if the file starts with resource but it isn't a test
		if strings.HasPrefix(v.Name(), "resource_") && strings.HasSuffix(v.Name(), ".go") && !strings.HasSuffix(v.Name(), "_test.go") {
			tffiles[v.Name()] = filepath.Join(internalPath, v.Name())
		}
	}

	// Check if we at least found one file
	if len(tffiles) == 0 {
		errexit("No terraform definition files have been found. You need a %q file as well as multiple %q files to have your provider run.", "provider.go", "resource_*")
	}

	// Check if we have a provider
	if _, found := tffiles["provider.go"]; !found {
		errexit("No %q found. This file is needed to provide the initial provider definition to your code.", "provider.go")
	}

	// Now check if once we have the provider, we have at least some other file
	if len(tffiles) == 1 {
		errexit("No resource definition files found. Resource files are needed to define a minimal documentation for the provider.")
	}

	// Create a token fileset
	fset := token.NewFileSet()

	// Iterate over each one of the files
	results := make([]parsedData, 0, len(tffiles))
	for _, v := range tffiles {
		results = append(results, parseProvider(internalName, fset, v))
	}

	// Print the output
	for _, v := range results {
		fmt.Println("Filename:", v.Filename)
		fmt.Println("Is provider config?:", v.IsMainProvider)
		fmt.Println("Resource name:", v.ResourceName)

		for _, j := range v.ResourceDefinition {
			fmt.Println("\t ", fmt.Sprintf("%#v", j))
		}
	}

}
