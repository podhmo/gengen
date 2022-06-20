//go:build gen

package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"go/types"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/podhmo/flagstruct"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

//go:embed gen.go.tmpl
var code string

type Input struct {
	Name     string
	Packages []string
	Schemas  []Symbol
}

type Symbol struct {
	Prefix string
	Name   string
}

type Options struct {
	PackageName string   `flag:"package-name"`
	Packages    []string `flag:"packages"`
	Interface   string   `flag:"interface"`

	DumpOnly bool `flag:"dump-only"`
	// TODO: debug option
}

func main() {
	options := &Options{PackageName: "gen", Interface: "github.com/podhmo/gengen/enum/schema.Interface"}
	if err := flagstruct.Parse(options); err != nil {
		log.Fatalf("parse options: %v", err)
	}

	if err := run(*options); err != nil {
		log.Fatalf("!! %+v", err)
	}
}

func run(options Options) error {
	var iface *types.Interface
	if !strings.Contains(options.Interface, ".") {
		return fmt.Errorf("--interface with <package>.<name>")
	}
	parts := strings.Split(options.Interface, ".")
	pathAndName := []string{strings.Join(parts[:len(parts)-1], "."), parts[len(parts)-1]}

	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedModule,
		Dir:  ".",
	}, append([]string{pathAndName[0]}, options.Packages...)...)

	if err != nil {
		return fmt.Errorf("package load: %w", err)
	}
	if len(pkgs) == 0 {
		return fmt.Errorf("no packages found")
	}

	targetPackages := make([]string, len(pkgs))
	for i, pkg := range pkgs {
		targetPackages[i] = pkg.PkgPath
	}

	for _, pkg := range pkgs {
		if pkg.PkgPath != pathAndName[0] {
			continue
		}

		ob := pkg.Types.Scope().Lookup(pathAndName[1])
		if ob == nil {
			return fmt.Errorf("%q is not found in %q", pathAndName[1], pkg.PkgPath)
		}

		t, ok := ob.Type().Underlying().(*types.Interface) // *Named -> *Interface
		if !ok {
			return fmt.Errorf("%q is not interface (type=%T)", options.Interface, ob.Type().Underlying())
		}
		iface = t
		break

	}
	if iface == nil {
		return fmt.Errorf("interface %q is not found", options.Interface)
	}

	var schemas []Symbol
	for _, pkg := range pkgs {
		prefix := pkg.Name
		// TODO: recursive
		s := pkg.Types.Scope()
		if pkg.PkgPath == pathAndName[0] {
			continue
		}

		for _, name := range s.Names() {
			ob := s.Lookup(name)
			typ, ok := ob.Type().(*types.Named)
			if !ok || !ob.Exported() {
				continue
			}
			if !types.Implements(typ, iface) {
				fmt.Println("\t\t not implements", typ)
				continue
			}
			schemas = append(schemas, Symbol{Prefix: prefix, Name: ob.Name()})
		}
	}

	data := Input{
		Name:     options.PackageName,
		Packages: targetPackages,
		Schemas:  schemas,
	}

	if options.DumpOnly {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(data)
	}

	var buf bytes.Buffer
	t := template.Must(template.New("gen-main").Parse(code))
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}
	b, err := imports.Process("main.go", buf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("format: %w", err)
	}

	fmt.Println(string(b))
	return nil
}
