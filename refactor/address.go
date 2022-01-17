package refactor

import "strings"

type ElementType string

const (
	Resource ElementType = "resource"
	Output               = "output"
	Var                  = "var"
	Local                = "local"
	Data                 = "data"
	Module               = "module"
)

type Address struct {
	elementType ElementType
	labels      []string
}

func (a *Address) RefNameArray() []string {
	if a.elementType == Resource {
		return a.labels
	}
	return append([]string{string(a.elementType)}, a.labels...)
}

func (a *Address) RefName() string {
	return strings.Join(a.RefNameArray(), ".")
}

func ParseAddress(addr string) *Address {
	parts := strings.Split(addr, ".")
	switch parts[0] {
	case string(Resource), "resources":
		return &Address{
			elementType: Resource,
			labels:      parts[1:],
		}
	case string(Output), "outputs":
		return &Address{
			elementType: Output,
			labels:      parts[1:],
		}
	case string(Var), "vars", "variable", "variables":
		return &Address{
			elementType: Output,
			labels:      parts[1:],
		}
	case string(Local), "locals":
		return &Address{
			elementType: Local,
			labels:      parts[1:],
		}
	case string(Data):
		return &Address{
			elementType: Data,
			labels:      parts[1:],
		}
	case string(Module):
		return &Address{
			elementType: Module,
			labels:      parts[1:],
		}
	default:
		return &Address{
			elementType: Resource,
			labels:      parts,
		}
	}
}
