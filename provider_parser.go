package main

import (
	"encoding/json"
	"fmt"
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

type docDefinition struct {
	Name     string
	Type     string
	Optional bool
	Computed bool
	Required bool
	ForceNew bool
	Inner    []docDefinition
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

	// Call the walker by creating a reference first
	// so we can get the list of attributes after
	wk := NewWalker(parsed.Decls)
	wk.Do()

	// Generate a structure with the proper details
	return parsedData{
		Filename:           filename,
		ResourceName:       resourcename,
		IsMainProvider:     resourcename == providerName,
		ResourceDefinition: wk.results,
	}
}

type mywalker struct {
	items   interface{}
	results []docDefinition
}

func NewWalker(nodeList interface{}) *mywalker {
	return &mywalker{items: nodeList}
}

func (w *mywalker) Do() {
	// Find if it's the right type
	list, ok := w.items.([]ast.Decl)
	if !ok {
		return
	}

	for _, v := range list {
		// Find if it's a function declaration
		fd, ok := v.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// Find if the body list contains at least
		// one element
		if len(fd.Body.List) != 1 {
			continue
		}

		// Find if that single element is a return
		// statement
		rs, ok := fd.Body.List[0].(*ast.ReturnStmt)
		if !ok {
			continue
		}

		// The element we're looking for has only
		// 1 return statement
		if len(rs.Results) != 1 {
			continue
		}

		// Type assert the value to be an UnaryExpr
		ue, ok := rs.Results[0].(*ast.UnaryExpr)
		if !ok {
			continue
		}

		// Find if this first element is a CompositeLit
		cl, ok := ue.X.(*ast.CompositeLit)
		if !ok {
			continue
		}

		// Holder for the element who holds the schema
		var kve *ast.KeyValueExpr

		// Find the schema in the list of elements
		for _, v := range cl.Elts {
			switch m := v.(type) {
			case *ast.KeyValueExpr:
				if m.Key.(*ast.Ident).Name == "Schema" {
					kve = m
				}
			}
		}

		// Check if we found something
		if kve == nil {
			continue
		}

		// Now the value for this thing is a CompositeLit too
		clt, ok := kve.Value.(*ast.CompositeLit)
		if !ok {
			continue
		}

		// We're looking for a map with all the instructions
		if _, ok := clt.Type.(*ast.MapType); !ok {
			continue
		}

		// Call the inner walker, we got the outer one
		iw := NewWalker(clt.Elts)
		iw.innerWalker()

		w.results = iw.results
	}
}

func (w *mywalker) innerWalker() {
	// Find if it's the right type
	list, ok := w.items.([]ast.Expr)
	if !ok {
		return
	}

	// Clean the list
	w.results = nil

	for _, v := range list {
		// Since we got called from the outside, we know
		// this is the definition list, so we expect only
		// an *ast.KeyValueExpr
		kve, ok := v.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		// Find if the object is a BasicLit, which is what we're expecting
		bl, ok := kve.Key.(*ast.BasicLit)
		if !ok {
			continue
		}

		// ... And the second parameter is always a CompositeLit which
		// defines the attributes of each schema element...
		val, ok := kve.Value.(*ast.CompositeLit)
		if !ok {
			continue
		}

		// Holder for all the definitions of the given object
		dd := docDefinition{
			Name: strings.Trim(bl.Value, `"`),
		}

		// Inside the values, we will find the details of all the schema definition
		// flags, such as mandatory, type, computed or optional

		for _, v := range val.Elts {
			// But to do so, the first one is always a KeyValueExpr
			e, ok := v.(*ast.KeyValueExpr)
			if !ok {
				continue
			}

			// And the key of the second value, which holds the schema requirements
			// is always an Ident, where the name holds the key
			kvk, ok := e.Key.(*ast.Ident)
			if !ok {
				continue
			}

			// Simple helper function
			getTrue := func(boolean *ast.KeyValueExpr) bool {
				return boolean.Value.(*ast.Ident).Name == "true"
			}

			// Depending on the value's key name, we can find what constraint
			// we're talking about
			switch kvk.Name {

			// Type is easy, just get the type, and remove the Prefix "Type***"
			case "Type":
				dd.Type = strings.TrimPrefix(e.Value.(*ast.SelectorExpr).Sel.Name, "Type")
				dd.Type = strings.ToLower(dd.Type)

			// For all the other ones, we need to cast the type to an actual ident
			// then get the Ident's name.
			case "Optional":
				dd.Optional = getTrue(e)
			case "ForceNew":
				dd.ForceNew = getTrue(e)
			case "Computed":
				dd.Computed = getTrue(e)
			case "Required":
				dd.Required = getTrue(e)
			}

			// Check now if type list, then we need to dereference all the inner parameters
			if (dd.Type == "list" || dd.Type == "set") && kvk.Name == "Elem" {
				// Get the inner value. If we pass the previous validations we can
				// assume it's a terraform project. Otherwise this double type assertion
				// will fail
				composite := e.Value.(*ast.UnaryExpr).X.(*ast.CompositeLit)

				// Get the kind, if this is a different Go project and this doesn't compile
				// then this part will fail with a panic, due to the 3 different type assertions
				kind := composite.Type.(*ast.SelectorExpr).Sel.Name

				// Switch depending on the kind
				if kind == "Schema" && len(composite.Elts) == 1 {
					// Generate a type based on this long type assertion
					// which if validated against a Terraform project, it's always
					// like this
					local := composite.Elts[0].(*ast.KeyValueExpr).Value.(*ast.SelectorExpr).Sel.Name

					// Then change the name to be "List (of []string)" or whatever the
					// actual type detected is
					if dd.Type == "list" {
						dd.Type = fmt.Sprintf("%s (of []%s)", dd.Type, strings.ToLower(strings.TrimPrefix(local, "Type")))
					} else {
						dd.Type = fmt.Sprintf("%s (of %s)", dd.Type, strings.ToLower(strings.TrimPrefix(local, "Type")))
					}
				}

				if kind == "Resource" {
					// Find the inner element, if this is a resource, it'll
					// have one inner element
					if len(composite.Elts) != 1 {
						continue
					}

					// The inner object of a resource is a KeyValueExpr
					resource, ok := composite.Elts[0].(*ast.KeyValueExpr)
					if !ok {
						continue
					}

					// Then, inside this resource we have an inner CompositeLit
					icl, ok := resource.Value.(*ast.CompositeLit)
					if !ok {
						continue
					}

					// Now walk inside these
					walk := NewWalker(icl.Elts)
					walk.innerWalker()

					// Once the walk is finished, grab those results and bring them back here
					dd.Inner = walk.results
					dd.Type = fmt.Sprintf("%s (of []object)", dd.Type)
				}
			}
		}

		// Once we've found all the different constraints, add them to the list of results
		w.results = append(w.results, dd)
	}
}

func pp(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return "----\n" + string(b) + "\n----\n\n"
}
