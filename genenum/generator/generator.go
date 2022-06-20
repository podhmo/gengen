package generator

import (
	"bytes"
	"embed"
	"fmt"
	"sort"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/pkg/errors"
	"golang.org/x/tools/imports"
)

//go:embed enum.tmpl
var content embed.FS

// Generator is responsible for generating validation files for the given in a go source file.
type Generator struct {
	Version   string
	Revision  string
	BuildDate string
	BuiltBy   string

	t                 *template.Template
	knownTemplates    map[string]*template.Template
	userTemplateNames []string

	Config
}

type Config struct {
	NoPrefix        bool
	LowercaseLookup bool
	CaseInsensitive bool
	Marshal         bool
	Sql             bool
	Flag            bool
	Names           bool
	LeaveSnakeCase  bool
	Prefix          string
	SqlNullInt      bool
	SqlNullStr      bool
	Ptr             bool
	MustParse       bool
	ForceLower      bool
}

// NewGenerator is a constructor method for creating a new Generator with default
// templates loaded.
func NewGenerator() (*Generator, error) {

	funcs := sprig.TxtFuncMap()
	funcs["stringify"] = Stringify
	funcs["mapify"] = Mapify
	funcs["unmapify"] = Unmapify
	funcs["namify"] = Namify
	funcs["offset"] = Offset

	t, err := template.New("generator").Funcs(funcs).ParseFS(content, "enum.tmpl")
	if err != nil {
		return nil, err
	}

	g := &Generator{
		Version:           "-",
		Revision:          "-",
		BuildDate:         "-",
		BuiltBy:           "-",
		knownTemplates:    make(map[string]*template.Template),
		userTemplateNames: make([]string, 0),
		t:                 t,
		Config:            Config{},
	}

	g.updateTemplates()

	return g, nil
}

// Emit does the heavy lifting for the code generation starting from the parsed AST file.
func (g *Generator) Emit(pkg string, enums []Enum) ([]byte, error) {
	sort.Slice(enums, func(i, j int) bool { return enums[i].Name < enums[j].Name })

	vBuff := bytes.NewBuffer([]byte{})
	err := g.t.ExecuteTemplate(vBuff, "header", map[string]interface{}{
		"package":   pkg,
		"version":   g.Version,
		"revision":  g.Revision,
		"buildDate": g.BuildDate,
		"builtBy":   g.BuiltBy,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "Failed writing header")
	}

	var created int
	for name, enum := range enums {
		created++
		data := map[string]interface{}{
			"enum":       enum,
			"name":       name,
			"lowercase":  g.Config.LowercaseLookup,
			"nocase":     g.Config.CaseInsensitive,
			"marshal":    g.Config.Marshal,
			"sql":        g.Config.Sql,
			"flag":       g.Config.Flag,
			"names":      g.Config.Names,
			"ptr":        g.Config.Ptr,
			"sqlnullint": g.Config.SqlNullInt,
			"sqlnullstr": g.Config.SqlNullStr,
			"mustparse":  g.Config.MustParse,
			"forcelower": g.Config.ForceLower,
		}

		err = g.t.ExecuteTemplate(vBuff, "enum", data)
		if err != nil {
			return vBuff.Bytes(), errors.WithMessage(err, fmt.Sprintf("Failed writing enum data for enum: %q", name))
		}

		for _, userTemplateName := range g.userTemplateNames {
			err = g.t.ExecuteTemplate(vBuff, userTemplateName, data)
			if err != nil {
				return vBuff.Bytes(), errors.WithMessage(err, fmt.Sprintf("Failed writing enum data for enum: %q, template: %v", name, userTemplateName))
			}
		}
	}

	if created < 1 {
		// Don't save anything if we didn't actually generate any successful enums.
		return nil, nil
	}

	formatted, err := imports.Process(pkg, vBuff.Bytes(), nil)
	if err != nil {
		err = fmt.Errorf("generate: error formatting code %s\n\n%s", err, vBuff.String())
	}
	return formatted, err
}

// updateTemplates will update the lookup map for validation checks that are
// allowed within the template engine.
func (g *Generator) updateTemplates() {
	for _, template := range g.t.Templates() {
		g.knownTemplates[template.Name()] = template
	}
}
