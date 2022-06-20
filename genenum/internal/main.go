//go:build gen

package main

import (
	_ "embed"
	"log"
	"os"
	"text/template"
)

//go:embed gen.go.tmpl
var code string

type Input struct {
	Name    string
	Pkgs    []string
	Schemas []Symbol
}

type Symbol struct {
	Prefix string
	Name   string
}

func main() {
	t := template.Must(template.New("gen-main").Parse(code))
	data := Input{
		Name: "gen",
		Pkgs: []string{"github.com/podhmo/gengen/genenum/sandbox"},
		Schemas: []Symbol{
			{Prefix: "sandbox", Name: "Op"},
		},
	}

	err := t.Execute(os.Stdout, data)
	if err != nil {
		log.Fatalf("!! %+v", err)
	}
}
