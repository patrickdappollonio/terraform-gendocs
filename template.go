package main

import (
	"bytes"
	"fmt"
	"io"
	"unicode/utf8"
)

const (
	itemHeading = `%s "%s" {`
	innerItem   = `%s # type: %s, required: %v, optional: %v, computed: %v, force-new: %v, child: %v`
)

func printOutput(w io.Writer, data []parsedData) {
	b := &bytes.Buffer{}

	for _, v := range data {
		preheading := "resource"
		if v.IsMainProvider {
			preheading = "provider"
		}

		// Print the initial part like `provider "myprovider" {`
		fmt.Fprintf(b, itemHeading, preheading, v.ResourceName)
		fmt.Fprintln(b, "")

		printChildList(1, b, v.ResourceDefinition)

		// Close the bracket and line break
		fmt.Fprintln(b, "}")
		fmt.Fprintln(b, "")
	}

	fmt.Fprintln(w, b.String())
}

func printChildList(indentspaces int, w io.Writer, items []docDefinition) {
	length := 0
	for _, v := range items {
		if newlength := utf8.RuneCountInString(v.Name); length < newlength {
			length = newlength
		}
	}

	spaces := ind(indentspaces)
	indentspaces += 1
	for _, v := range items {
		fmt.Fprintf(w, spaces+innerItem, pad(v.Name, length), v.Type, v.Required, v.Optional, v.Computed, v.ForceNew, len(v.Inner))
		fmt.Fprintln(w, "")
		printChildList(indentspaces, w, v.Inner)
	}
}

func pad(str string, length int) string {
	return str + times(" ", length-utf8.RuneCountInString(str))
}

func ind(n int) string {
	return times("\t", n)
}

func times(str string, n int) (out string) {
	for i := 1; i <= n; i++ {
		out += str
	}

	return
}
