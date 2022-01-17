package refactor

import "strings"

type TypeName string

const (
	TypeResource TypeName = "resource"
	TypeOutput            = "output"
	TypeVar               = "var"
	TypeLocal             = "local"
	TypeData              = "data"
	TypeModule            = "module"
)

type Address struct {
	elementType TypeName
	labels      []string
}

func (a *Address) RefNameArray() []string {
	if a.elementType == TypeResource {
		return a.labels
	}
	return append([]string{string(a.elementType)}, a.labels...)
}

func (a *Address) RefName() string {
	return strings.Join(a.RefNameArray(), ".")
}

func (a *Address) BlockType() string {
	if a.elementType == TypeVar {
		return "variable"
	}
	return string(a.elementType)
}

func ParseAddress(addr string) *Address {
	parts := strings.Split(addr, ".")
	switch parts[0] {
	case string(TypeResource), "resources":
		return &Address{
			elementType: TypeResource,
			labels:      parts[1:],
		}
	case string(TypeOutput), "outputs":
		return &Address{
			elementType: TypeOutput,
			labels:      parts[1:],
		}
	case string(TypeVar), "vars", "variable", "variables":
		return &Address{
			elementType: TypeVar,
			labels:      parts[1:],
		}
	case string(TypeLocal), "locals":
		return &Address{
			elementType: TypeLocal,
			labels:      parts[1:],
		}
	case string(TypeData):
		return &Address{
			elementType: TypeData,
			labels:      parts[1:],
		}
	case string(TypeModule):
		return &Address{
			elementType: TypeModule,
			labels:      parts[1:],
		}
	default:
		return &Address{
			elementType: TypeResource,
			labels:      parts,
		}
	}
}
