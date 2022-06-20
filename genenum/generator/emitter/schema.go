package emitter

import (
	"encoding/json"
	"fmt"
)

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

func UnmarshalSchema(b []byte) (Enum, error) {
	var dst Enum
	if err := json.Unmarshal(b, &dst); err != nil {
		return dst, fmt.Errorf("failed to unmarshal schema: %w", err)
	}

	// fixme: toJSON uint as float64 in marshal
	for i, v := range dst.Values {
		v.Value = uint64(v.Value.(float64))
		dst.Values[i] = v
	}
	return dst, nil
}
