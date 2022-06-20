package schema

import (
	"encoding/json"
	"reflect"

	"github.com/podhmo/gengen/genenum/generator/load"
)

type Interface interface {
	EnumName(self interface{}) string
	EnumType() string
	EnumValues() []*EnumValue
}

type Enum struct {
	c counter
}

func (e Enum) EnumName(self interface{}) string {
	return reflect.TypeOf(self).Name()
}

func (e Enum) EnumType() string {
	return "uint64" // TODO: implementation
}

func (e Enum) EnumValues() []*EnumValue {
	return nil
}

func (e *Enum) Uint(name string) *EnumValue {
	value := e.c.Count()
	return &EnumValue{value: uint64(value), name: name}
}

type EnumValue struct {
	name    string
	value   interface{}
	comment string
}

func (v *EnumValue) Name(name string) *EnumValue {
	v.name = name
	return v
}

func (v *EnumValue) Comment(value string) *EnumValue {
	v.comment = value
	return v
}

func toLoadSchema(e Interface) load.Enum {
	src := e.EnumValues()
	dst := make([]load.EnumValue, len(src))
	for i, x := range src {
		dst[i] = load.EnumValue{
			RawName:      x.name,
			Name:         x.name,
			PrefixedName: x.name,
			Value:        x.value,
			Comment:      x.comment,
		}
	}

	return load.Enum{
		Name:   e.EnumName(e),
		Type:   e.EnumType(),
		Prefix: "",
		Values: dst,
	}
}

func MarshalSchema(e Interface) ([]byte, error) {
	return json.Marshal(toLoadSchema(e))
}

// util
type counter struct {
	c uint
}

func (c *counter) Reset(z uint) {
	c.c = z
}

func (c *counter) Count() uint {
	v := c.c
	c.c++
	return v
}
