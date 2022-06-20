package enum

import (
	"encoding/json"
	"reflect"

	"github.com/podhmo/gengen/enum/generator/load"
)

type Interface interface {
	EnumName(self interface{}) string
	EnumValues() []*Value
}

type Schema struct {
	c counter
}

func (e Schema) EnumName(self interface{}) string {
	return reflect.TypeOf(self).Name()
}

func (e Schema) EnumValues() []*Value {
	return nil
}

func (e *Schema) Uint(name string) *Value {
	value := e.c.Count()
	return &Value{value: uint64(value), name: name}
}

type Value struct {
	name    string
	value   interface{}
	comment string
}

func (v *Value) Name(name string) *Value {
	v.name = name
	return v
}

func (v *Value) Comment(value string) *Value {
	v.comment = value
	return v
}

func guessType(src []*Value) string {
	for _, x := range src {
		return reflect.TypeOf(x.value).Name()
	}
	return "string" // ?
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

	typ := guessType(src)
	return load.Enum{
		Name:   e.EnumName(e),
		Prefix: "",
		Type:   typ,
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
