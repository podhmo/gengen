package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/podhmo/gengen/enum"
	"github.com/podhmo/gengen/enum/internal/sandbox"
)

type Package struct {
	Name string
	Defs []json.RawMessage
}

// pkg name
var pkg = "gen"

// input:
var schemas = []enum.Interface{
	sandbox.Op{},
}

func run() error {
	pkg := Package{Name: pkg}
	for _, x := range schemas {
		typename := x.EnumName(x)

		b, err := enum.MarshalSchema(x)
		if err != nil {
			return fmt.Errorf("marshal in %v: %w", typename, err)
		}
		pkg.Defs = append(pkg.Defs, b)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(pkg)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "!! %+v", err)
		os.Exit(1)
	}
}
