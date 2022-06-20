package schema

import "reflect"

// TODO(podhmo): to interface

type Counter struct {
	c uint
}

func (c *Counter) Reset(z uint) {
	c.c = z
}

func (c *Counter) Count() uint {
	v := c.c
	c.c++
	return v
}

type EnumInterface interface {
	Name(self interface{}) string
	Type() string
	Values() []*EnumValue
}

type Enum struct {
	c Counter
}

func (e *Enum) Name(self interface{}) string {
	return reflect.TypeOf(self).Elem().Name()
}

func (e *Enum) Type() string {
	return "uint64" // TODO: implementation
}

func (e *Enum) Values() []*EnumValue {
	return nil
}

func (e *Enum) Uint(name string) *EnumValue {
	value := e.c.Count()
	return &EnumValue{value: uint64(value), name: name}
}

type EnumValue struct {
	name  string
	value interface{}
}

func (v *EnumValue) Name(name string) *EnumValue {
	v.name = name
	return v
}

// ToEmitterInput : emitter.Enum
func ToEmitterInput(e EnumInterface) map[string]interface{} {
	src := e.Values()
	dst := make([]map[string]interface{}, len(src))
	for i, x := range src {
		dst[i] = map[string]interface{}{
			"RawName":      x.name,
			"Name":         x.name,
			"PrefixedName": x.name,
			"Value":        x.value,
			"Comment":      "",
		}
	}

	return map[string]interface{}{
		"Name":   e.Name(e),
		"Type":   e.Type(),
		"Prefix": "-",
		"Values": dst,
	}
}
