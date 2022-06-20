package internal

import (
	"github.com/podhmo/gengen/genenum/schema"
)

type Op struct {
	schema.Enum
}

func (op Op) EnumValues() []*schema.EnumValue {
	return []*schema.EnumValue{
		op.Uint("Add").Comment("lhs + rhs"),
		op.Uint("Mul"),
		op.Uint("Sub"),
	}
}
