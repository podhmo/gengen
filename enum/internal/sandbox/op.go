package sandbox

import (
	"github.com/podhmo/gengen/enum"
)

type Op struct {
	enum.Schema
}

func (op Op) EnumValues() []*enum.Value {
	return []*enum.Value{
		op.Uint("Add").Comment("lhs + rhs"),
		op.Uint("Mul"),
		op.Uint("Sub"),
	}
}
