package emitter

import (
	"bytes"
	"embed"
	"fmt"
	"go/token"
	"sort"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/pkg/errors"
	"golang.org/x/tools/imports"
)

//go:embed enum.tmpl
var content embed.FS

// Emitter is responsible for generating validation files for the given in a go source file.
type Emitter struct {
	Version           string
	Revision          string
	BuildDate         string
	BuiltBy           string
	t                 *template.Template
	knownTemplates    map[string]*template.Template
	userTemplateNames []string
	fileSet           *token.FileSet
	noPrefix          bool
	lowercaseLookup   bool
	caseInsensitive   bool
	marshal           bool
	sql               bool
	flag              bool
	names             bool
	leaveSnakeCase    bool
	prefix            string
	sqlNullInt        bool
	sqlNullStr        bool
	ptr               bool
	mustParse         bool
	forceLower        bool
}

// Enum holds data for a discovered enum in the parsed source
type Enum struct {
	Name   string
	Prefix string
	Type   string
	Values []EnumValue
}

// EnumValue holds the individual data for each enum value within the found enum.
type EnumValue struct {
	RawName      string
	Name         string
	PrefixedName string
	Value        interface{}
	Comment      string
}

// NewEmitter is a constructor method for creating a new Emitter with default
// templates loaded.
func NewEmitter() (*Emitter, error) {

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

	g := &Emitter{
		Version:           "-",
		Revision:          "-",
		BuildDate:         "-",
		BuiltBy:           "-",
		knownTemplates:    make(map[string]*template.Template),
		userTemplateNames: make([]string, 0),
		t:                 t,
		fileSet:           token.NewFileSet(), // todo: remove
		noPrefix:          false,
	}

	g.updateTemplates()

	return g, nil
}

// Emit does the heavy lifting for the code generation starting from the parsed AST file.
func (g *Emitter) Emit(pkg string, enums []Enum) ([]byte, error) {
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
			"lowercase":  g.lowercaseLookup,
			"nocase":     g.caseInsensitive,
			"marshal":    g.marshal,
			"sql":        g.sql,
			"flag":       g.flag,
			"names":      g.names,
			"ptr":        g.ptr,
			"sqlnullint": g.sqlNullInt,
			"sqlnullstr": g.sqlNullStr,
			"mustparse":  g.mustParse,
			"forcelower": g.forceLower,
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
func (g *Emitter) updateTemplates() {
	for _, template := range g.t.Templates() {
		g.knownTemplates[template.Name()] = template
	}
}
