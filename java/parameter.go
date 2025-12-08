package java

type Parameter struct {
	Name  string
	Type  Type
	Index int
}

func (p Parameter) String() string {
	if p.Name != "" {
		return p.Type.String() + " " + p.Name
	}
	return p.Type.String()
}
