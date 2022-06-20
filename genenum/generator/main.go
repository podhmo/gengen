//go:buld gen

package main

import (
	"fmt"
	"log"

	"github.com/podhmo/gengen/genenum/generator/emitter"
)

func run() error {
	g, err := emitter.NewEmitter()
	if err != nil {
		return fmt.Errorf("new emitter: %w", err)
	}

	enums := []emitter.Enum{
		{
			Name: "Op",
			Type: "uint8",
			Values: []emitter.EnumValue{
				{RawName: "Add", PrefixedName: "Add", Name: "Add", Value: uint64(1)}, // RawName, PrefixedName, Comment
				{RawName: "Sub", PrefixedName: "Sub", Name: "Sub", Value: uint64(2)},
				{RawName: "Mul", PrefixedName: "Mul", Name: "Mul", Value: uint64(3)},
			}},
	}
	pkg := "gen"
	b, err := g.Emit(pkg, enums)
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
