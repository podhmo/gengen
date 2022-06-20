package schema

import (
	"encoding/json"
	"reflect"
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

func toGeneratorInput(e Interface) map[string]interface{} {
	src := e.EnumValues()
	dst := make([]map[string]interface{}, len(src))
	for i, x := range src {
		dst[i] = map[string]interface{}{
			"RawName":      x.name,
			"Name":         x.name,
			"PrefixedName": x.name,
			"Value":        x.value,
			"Comment":      x.comment,
		}
	}

	return map[string]interface{}{
		"Name":   e.EnumName(e),
		"Type":   e.EnumType(),
		"Prefix": "-",
		"Values": dst,
	}
}

func MarshalSchema(e Interface) ([]byte, error) {
	return json.Marshal(toGeneratorInput(e))
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
