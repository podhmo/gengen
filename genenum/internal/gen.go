//go:buld gen

package main

import (
	"encoding/json"
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
	schemas := []schema.EnumInterface{
		sandbox.Op{},
	}

	targets := make([]emitter.Enum, len(schemas))

	for i, x := range schemas {
		var dst emitter.Enum
		typename := x.Name(x)
		src := schema.ToEmitterInput(x)

		b, err := json.Marshal(src)
		if err != nil {
			return fmt.Errorf("marshal in %v: %w", typename, err)
		}

		if err := json.Unmarshal(b, &dst); err != nil {
			return fmt.Errorf("unmarshal in %v: %w", typename, err)
		}

		// fixme: toJSON uint as float64 in marshal
		for i, v := range dst.Values {
			v.Value = uint64(v.Value.(float64))
			dst.Values[i] = v
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
