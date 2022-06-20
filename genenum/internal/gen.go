//go:buld gen

// Code generated by gen.go.tmpl DO NOT EDIT.
package main

import (
	"fmt"
	"log"

	"github.com/podhmo/gengen/genenum/generator/emitter"
	"github.com/podhmo/gengen/genenum/schema"

	"github.com/podhmo/gengen/genenum/sandbox"
)

func run() error {
	g, err := emitter.NewEmitter()
	if err != nil {
		return fmt.Errorf("new emitter: %w", err)
	}

	// pkg name
	pkg := "gen"
	// input:
	schemas := []schema.Interface{
		sandbox.Op{},
	}

	targets := make([]emitter.Enum, len(schemas))

	for i, x := range schemas {
		typename := x.EnumName(x)

		b, err := schema.MarshalSchema(x)
		if err != nil {
			return fmt.Errorf("marshal in %v: %w", typename, err)
		}
		dst, err := emitter.UnmarshalSchema(b)
		if err != nil {
			return fmt.Errorf("unmarshal in %v: %w", typename, err)
		}

		targets[i] = dst
	}

	b, err := g.Emit(pkg, targets)
	if err != nil {
		return fmt.Errorf("emit: %w", err)
	}
	fmt.Println(string(b))
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("!! %+v", err)
	}
}

