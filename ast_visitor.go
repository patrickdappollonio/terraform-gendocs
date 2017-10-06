package main

import (
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
	switch t := node.(type) {

	// The type we're always looking which contains all the defined
	// data is a KeyValueExpr.
	case *ast.KeyValueExpr:
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
				if e.Value.(*ast.Ident).Name == "true" {
					dd.Optional = true
				}
			case "ForceNew":
				if e.Value.(*ast.Ident).Name == "true" {
					dd.ForceNew = true
				}
			case "Computed":
				if e.Value.(*ast.Ident).Name == "true" {
					dd.Computed = true
				}
			}
		}

		// Once we've found all the different constraints, add them to the list of results
		w.results = append(w.results, dd)
	}

	return w
}
