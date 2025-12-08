package classfile

import "strings"

type FieldType struct {
	BaseType   string
	ClassName  string
	ArrayDepth int
}

func (ft *FieldType) String() string {
	var sb strings.Builder
	for i := 0; i < ft.ArrayDepth; i++ {
		sb.WriteString("[]")
	}
	if ft.BaseType != "" {
		sb.WriteString(ft.BaseType)
	} else if ft.ClassName != "" {
		sb.WriteString(strings.ReplaceAll(ft.ClassName, "/", "."))
	}
	return sb.String()
}

func (ft *FieldType) IsArray() bool {
	return ft.ArrayDepth > 0
}

func (ft *FieldType) IsPrimitive() bool {
	return ft.BaseType != "" && ft.ClassName == ""
}

func (ft *FieldType) IsReference() bool {
	return ft.ClassName != "" || ft.ArrayDepth > 0
}

type MethodDescriptor struct {
	Parameters []FieldType
	ReturnType *FieldType
}

func (md *MethodDescriptor) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	for i, p := range md.Parameters {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(p.String())
	}
	sb.WriteString(")")
	if md.ReturnType != nil {
		sb.WriteString(" ")
		sb.WriteString(md.ReturnType.String())
	} else {
		sb.WriteString(" void")
	}
	return sb.String()
}

func ParseFieldDescriptor(desc string) *FieldType {
	ft, _ := parseFieldType(desc, 0)
	return ft
}

func ParseMethodDescriptor(desc string) *MethodDescriptor {
	if len(desc) == 0 || desc[0] != '(' {
		return nil
	}

	md := &MethodDescriptor{}
	i := 1

	for i < len(desc) && desc[i] != ')' {
		ft, consumed := parseFieldType(desc, i)
		if ft == nil {
			return nil
		}
		md.Parameters = append(md.Parameters, *ft)
		i += consumed
	}

	if i >= len(desc) || desc[i] != ')' {
		return nil
	}
	i++

	if i < len(desc) {
		if desc[i] == 'V' {
			md.ReturnType = nil
		} else {
			md.ReturnType, _ = parseFieldType(desc, i)
		}
	}

	return md
}

func parseFieldType(desc string, start int) (*FieldType, int) {
	if start >= len(desc) {
		return nil, 0
	}

	ft := &FieldType{}
	i := start

	for i < len(desc) && desc[i] == '[' {
		ft.ArrayDepth++
		i++
	}

	if i >= len(desc) {
		return nil, 0
	}

	switch desc[i] {
	case 'B':
		ft.BaseType = "byte"
		return ft, i - start + 1
	case 'C':
		ft.BaseType = "char"
		return ft, i - start + 1
	case 'D':
		ft.BaseType = "double"
		return ft, i - start + 1
	case 'F':
		ft.BaseType = "float"
		return ft, i - start + 1
	case 'I':
		ft.BaseType = "int"
		return ft, i - start + 1
	case 'J':
		ft.BaseType = "long"
		return ft, i - start + 1
	case 'S':
		ft.BaseType = "short"
		return ft, i - start + 1
	case 'Z':
		ft.BaseType = "boolean"
		return ft, i - start + 1
	case 'L':
		semicolon := strings.IndexByte(desc[i:], ';')
		if semicolon == -1 {
			return nil, 0
		}
		ft.ClassName = desc[i+1 : i+semicolon]
		return ft, i - start + semicolon + 1
	default:
		return nil, 0
	}
}

func InternalToSourceName(name string) string {
	return strings.ReplaceAll(name, "/", ".")
}

func SourceToInternalName(name string) string {
	return strings.ReplaceAll(name, ".", "/")
}
