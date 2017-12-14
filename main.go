package main

import (
	"fmt"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//go:generate go run generator.go

const help = "Usage: terraform-gendocs $IMPORT_PATH $FORMAT{html|hcl} [$FILENAME]"

func main() {
	// Check if it's help
	if len(os.Args) == 2 {
		if os.Args[1] == "--help" || os.Args[1] == "-h" {
			fmt.Fprintln(os.Stderr, help)
			os.Exit(0)
		}
	}

	// Find if a second argument is passed
	if l := len(os.Args); l < 3 || l > 4 {
		errexit(help)
	}

	// Define the format
	exportFormat := os.Args[2]
	if exportFormat != "html" && exportFormat != "hcl" {
		errexit("Wrong 'format' parameter. Use 'html' or 'hcl'.\n%s", help)
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
		if pp := parseProvider(internalName, fset, v); pp.ResourceName != "" {
			results = append(results, pp)
		}
	}

	// Print the output
	sort.Sort(byProviderResource(results))

	// Output filename
	var (
		fileName = fmt.Sprintf("%s-docs", internalName)
		attrib   = int(os.O_RDWR | os.O_CREATE)
		perms    = os.FileMode(0644)
	)

	// Check what we will use
	if len(os.Args) == 4 {
		if dn := slugify(os.Args[3]); dn != "" {
			fileName = dn
		}
	}

	switch exportFormat {
	case "hcl":
		// Create an exported file (or update if it already exists) with proper permissions
		f, err := os.OpenFile(fmt.Sprintf("./%s.tf", fileName), attrib, perms)
		if err != nil {
			errexit("Unable to save HCL output: %s", err.Error())
		}

		// Close once we're ready
		defer f.Close()

		printOutput(f, results)
	case "html":
		// Create an exported file (or update if it already exists) with proper permissions
		f, err := os.OpenFile(fmt.Sprintf("./%s.html", fileName), attrib, perms)
		if err != nil {
			errexit("Unable to save HCL output: %s", err.Error())
		}

		// Close once we're ready
		defer f.Close()

		printHTMLOutput(f, results)
	}
}
