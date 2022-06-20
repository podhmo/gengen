package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/podhmo/flagstruct"
	"github.com/podhmo/gengen/enum/generator"
	"github.com/podhmo/gengen/enum/generator/load"
)

type Package struct {
	Name string
	Defs []load.Enum
}

type Options struct {
	InputFile string `flag:"input-file"`
}

func run(options Options) error {
	g, err := generator.NewGenerator()
	if err != nil {
		return fmt.Errorf("new generator: %w", err)
	}

	f, err := os.Open(options.InputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	var pkg Package
	if err := json.NewDecoder(f).Decode(&pkg); err != nil {
		return err
	}

	// fixme: toJSON uint as float64 in marshal
	for _, x := range pkg.Defs {
		for i, v := range x.Values {
			v.Value = uint64(v.Value.(float64))
			x.Values[i] = v
		}
	}

	fmt.Println("----------------------------------------")
	fmt.Println("package:", pkg.Name)
	fmt.Println("----------------------------------------")

	b, err := g.Emit(pkg.Name, pkg.Defs)
	if err != nil {
		return fmt.Errorf("emit: %w", err)
	}
	fmt.Println(string(b))
	return nil
}

func main() {
	options := &Options{}
	if err := flagstruct.Parse(options); err != nil {
		log.Fatalf("parse options %+v", err)
	}
	if err := run(*options); err != nil {
		log.Fatalf("!! %+v", err)
	}
}
