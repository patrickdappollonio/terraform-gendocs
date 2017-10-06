package main

import (
	"fmt"
	"go/ast"
	"strings"
)

type docDefinition struct {
	Name     string
	Type     string
	Optional bool
	Computed bool
	Required bool
	ForceNew bool
	Inner    []docDefinition
}

type visitorProvider struct {
	results []docDefinition
}

func (w *visitorProvider) Visit(node ast.Node) ast.Visitor {
	// Depending on the node type we do some operations
	t, ok := node.(*ast.KeyValueExpr)
	if !ok {
		return w
	}

	// ... Where the first parameter needs to be a BasicLit, which
	// contains the schema name...
	bl, ok := t.Key.(*ast.BasicLit)
	if !ok {
		return w
	}

	// ... And the second parameter is always a CompositeLit which
	// defines the attributes of each schema element...
	val, ok := t.Value.(*ast.CompositeLit)
	if !ok {
		return w
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
		if dd.Type == "list" && kvk.Name == "Elem" {
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
				dd.Type = fmt.Sprintf("%s (of []%s)", dd.Type, strings.ToLower(strings.TrimPrefix(local, "Type")))
			}

			if kind == "Resource" {
				// If it's a resource, we can actually use the same
				// visitor to fetch the internal details recursively
				internal := new(visitorProvider)
				ast.Walk(internal, composite)

				// Once the walk is finished, grab those results and bring them back here
				dd.Inner = internal.results
				dd.Type = fmt.Sprintf("%s (of []object)", dd.Type)
			}
		}
	}

	// Once we've found all the different constraints, add them to the list of results
	w.results = append(w.results, dd)

	return w
}
