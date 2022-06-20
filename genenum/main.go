//go:buld gen

package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/podhmo/gengen/genenum/generator/emitter"
	"github.com/podhmo/gengen/genenum/schema"
)

type Op struct {
	schema.Enum
}

func (op *Op) Values() []*schema.EnumValue {
	return []*schema.EnumValue{
		op.Uint("Add"),
		op.Uint("Mul"),
		op.Uint("Sub"),
	}
}

func run() error {
	g, err := emitter.NewEmitter()
	if err != nil {
		return fmt.Errorf("new emitter: %w", err)
	}

	// TODO: this
	op := &Op{}
	typename := op.Name(op)

	src := schema.ToEmitterInput(op)
	var dst emitter.Enum

	{
		b, err := json.Marshal(src)
		if err != nil {
			return fmt.Errorf("marshal in %v: %w", typename, err)
		}

		if err := json.Unmarshal(b, &dst); err != nil {
			return fmt.Errorf("unmarshal in %v: %w", typename, err)
		}
	}

	// fixme
	for i, v := range dst.Values {
		v.Value = uint64(v.Value.(float64))
		dst.Values[i] = v
	}

	pkg := "gen"
	enums := []emitter.Enum{dst}
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
