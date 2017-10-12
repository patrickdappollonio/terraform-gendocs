package main

import (
	"fmt"
	"html/template"
	"io"
	"path"
	"strings"
)

const outputName = "output.tmpl"

var (
	fm = template.FuncMap{
		"resname": func(s string) string {
			m := ""
			for _, v := range strings.Split(s, "/") {
				m += fmt.Sprintf("%s / ", v)
			}

			return strings.Trim(m, "/ ")
		},
		"yesno": func(b bool) string {
			if b {
				return "yes"
			}
			return "no"
		},
		"coltype": func(s string) template.HTML {
			// Get the type first
			elems := strings.Split(s, " ")

			// Check if we got at least one
			if len(elems) == 0 {
				return template.HTML("")
			}

			typed := strings.TrimSpace(elems[0])
			output := ""

			switch typed {
			case "list", "map", "set":
				// If there's no "two elements", then return
				if len(elems) != 3 {
					return template.HTML(fmt.Sprintf(`<code>%s</code>`, typed))
				}

				// Get the second part
				valued := strings.NewReplacer(
					" ", "",
					")", "",
				).Replace(elems[len(elems)-1])

				output = fmt.Sprintf(`<code>%s</code>`, valued)
			default:
				output = fmt.Sprintf(`<code>%s</code>`, typed)
			}

			return template.HTML(output)

		},
	}

	t = template.Must(template.New(outputName).Funcs(fm).ParseFiles(outputName))
)

type res struct {
	Resource       string
	IsMainProvider bool
	IsChild        bool
	Parameters     []docDefinition
}

func printHTMLOutput(w io.Writer, data []parsedData) {
	res := recursivelyFlatten(data)

	bundle := map[string]interface{}{
		"Title":   fmt.Sprintf("Documentation for %q Provider", findProviderName(data)),
		"Details": res,
	}

	if err := t.Execute(w, bundle); err != nil {
		errexit("Unable to generate output file: %s", err.Error())
	}
}

func findProviderName(data []parsedData) string {
	for _, v := range data {
		if v.IsMainProvider {
			return v.ResourceName
		}
	}

	return ""
}

func recursivelyFlatten(data []parsedData) []res {
	var out []res

	for _, v := range data {
		out = append(out, res{
			Resource:       v.ResourceName,
			IsMainProvider: v.IsMainProvider,
			IsChild:        strings.Contains(v.ResourceName, "/"),
			Parameters:     v.ResourceDefinition,
		})

		for _, m := range v.ResourceDefinition {
			if len(m.Inner) > 0 {
				p := recursivelyFlatten([]parsedData{{
					ResourceName:       path.Join(v.ResourceName, m.Name),
					ResourceDefinition: m.Inner,
				}})

				out = append(out, p...)
			}
		}
	}

	return out
}
